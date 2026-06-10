from fastapi import FastAPI, Depends, HTTPException, status
from fastapi.security import HTTPBearer, HTTPAuthorizationCredentials
from sqlalchemy import create_engine, Column, Integer, String, Boolean
from sqlalchemy.ext.declarative import declarative_base
from sqlalchemy.orm import sessionmaker, Session
from jose import jwt, JWTError
from passlib.context import CryptContext
from pydantic import BaseModel
from typing import Optional, List

app = FastAPI()
security = HTTPBearer()
pwd_context = CryptContext(schemes=["bcrypt"], deprecated="auto")
SECRET_KEY = "secret"
ALGORITHM = "HS256"

# DB setup
SQLALCHEMY_DATABASE_URL = "sqlite:///./test.db"
engine = create_engine(SQLALCHEMY_DATABASE_URL)
SessionLocal = sessionmaker(autocommit=False, autoflush=False, bind=engine)
Base = declarative_base()

class User(Base):
    __tablename__ = "users"
    id = Column(Integer, primary_key=True, index=True)
    email = Column(String, unique=True, index=True)
    name = Column(String)
    hashed_password = Column(String)
    role = Column(String, default="user")

class Task(Base):
    __tablename__ = "tasks"
    id = Column(Integer, primary_key=True, index=True)
    title = Column(String, index=True)
    description = Column(String, nullable=True)
    status = Column(String, default="pending")
    assigned_to = Column(Integer, nullable=True)
    created_by = Column(Integer)

Base.metadata.create_all(bind=engine)

# Schemas
class UserCreate(BaseModel):
    email: str; name: str; password: str

class Token(BaseModel):
    access_token: str

# Helpers
def get_db():
    db = SessionLocal()
    try: yield db
    finally: db.close()

def get_current_user(credentials: HTTPAuthorizationCredentials = Depends(security), db: Session = Depends(get_db)):
    token = credentials.credentials
    try:
        payload = jwt.decode(token, SECRET_KEY, algorithms=[ALGORITHM])
        user_id = payload.get("sub")
        if user_id is None: raise HTTPException(401)
        user = db.query(User).filter(User.id == user_id).first()
        if user is None: raise HTTPException(401)
        return user
    except JWTError: raise HTTPException(401)

@app.post("/auth/register")
def register(user: UserCreate, db: Session = Depends(get_db)):
    hashed = pwd_context.hash(user.password)
    db_user = User(email=user.email, name=user.name, hashed_password=hashed)
    db.add(db_user); db.commit(); db.refresh(db_user)
    return db_user

@app.post("/auth/login", response_model=Token)
def login(email: str, password: str, db: Session = Depends(get_db)):
    user = db.query(User).filter(User.email == email).first()
    if not user or not pwd_context.verify(password, user.hashed_password):
        raise HTTPException(401)
    token = jwt.encode({"sub": str(user.id), "role": user.role}, SECRET_KEY, algorithm=ALGORITHM)
    return {"access_token": token}

@app.get("/tasks")
def list_tasks(page: int = 1, limit: int = 10, status: Optional[str] = None, current_user: User = Depends(get_current_user), db: Session = Depends(get_db)):
    query = db.query(Task)
    if current_user.role != "admin":
        query = query.filter(Task.assigned_to == current_user.id)
    if status:
        query = query.filter(Task.status == status)
    tasks = query.offset((page-1)*limit).limit(limit).all()
    return tasks
