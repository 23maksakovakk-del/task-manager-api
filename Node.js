
Ниже – **ключевые фрагменты кода** для каждой реализации (полный код даётся компактно, но достаточный для демонстрации).

---

### 1. Node.js (Express + Prisma + JWT)

**Основной файл `server.js`** (упрощённо, полная версия в предыдущем ответе)

```javascript
const express = require('express');
const jwt = require('jsonwebtoken');
const bcrypt = require('bcryptjs');
const { PrismaClient } = require('@prisma/client');
const prisma = new PrismaClient();

const app = express();
app.use(express.json());

// Middleware auth
const auth = (req, res, next) => {
  const token = req.header('Authorization')?.replace('Bearer ', '');
  try {
    req.user = jwt.verify(token, process.env.JWT_SECRET);
    next();
  } catch { res.status(401).json({ error: 'Unauthorized' }); }
};

// Регистрация
app.post('/auth/register', async (req, res) => {
  const { email, name, password } = req.body;
  const hashed = await bcrypt.hash(password, 10);
  const user = await prisma.user.create({ data: { email, name, password: hashed } });
  res.json(user);
});

// Логин
app.post('/auth/login', async (req, res) => {
  const user = await prisma.user.findUnique({ where: { email: req.body.email } });
  if (!user || !await bcrypt.compare(req.body.password, user.password))
    return res.status(401).json({ error: 'Invalid credentials' });
  const token = jwt.sign({ id: user.id, role: user.role }, process.env.JWT_SECRET);
  res.json({ token });
});

// GET /tasks с пагинацией и фильтром (только для user – свои задачи, admin – все)
app.get('/tasks', auth, async (req, res) => {
  const { page = 1, limit = 10, status } = req.query;
  const where = { status };
  if (req.user.role !== 'admin') where.assignedTo = req.user.id;
  const tasks = await prisma.task.findMany({
    where,
    skip: (page-1)*limit,
    take: +limit
  });
  res.json(tasks);
});

app.listen(3000);
