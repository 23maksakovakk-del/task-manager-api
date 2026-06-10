using Microsoft.AspNetCore.Authentication.JwtBearer;
using Microsoft.EntityFrameworkCore;
using Microsoft.IdentityModel.Tokens;
using System.Text;
using TaskManagerApi.Models;
using TaskManagerApi.Services;

var builder = WebApplication.CreateBuilder(args);
builder.Services.AddDbContext<AppDbContext>(opt => opt.UseNpgsql("Host=localhost;Database=taskdb;Username=postgres;Password=secret"));
builder.Services.AddAuthentication(JwtBearerDefaults.AuthenticationScheme).AddJwtBearer(opt =>
{
    opt.TokenValidationParameters = new TokenValidationParameters
    {
        ValidateIssuerSigningKey = true,
        IssuerSigningKey = new SymmetricSecurityKey(Encoding.UTF8.GetBytes("secretsecretsecret")),
        ValidateIssuer = false, ValidateAudience = false
    };
});
builder.Services.AddScoped<IAuthService, AuthService>();
builder.Services.AddScoped<ITaskService, TaskService>();
builder.Services.AddControllers();

var app = builder.Build();
app.UseAuthentication(); app.UseAuthorization();
app.MapControllers();
app.Run();

// AuthController.cs
[ApiController, Route("auth")]
public class AuthController : ControllerBase {
    [HttpPost("register")] public async Task<IActionResult> Register(RegisterDto dto, [FromServices] IAuthService auth) => Ok(await auth.Register(dto));
    [HttpPost("login")] public async Task<IActionResult> Login(LoginDto dto, [FromServices] IAuthService auth) => Ok(new { token = await auth.Login(dto) });
}

// TasksController.cs
[ApiController, Route("tasks"), Authorize]
public class TasksController : ControllerBase {
    [HttpGet] public async Task<IActionResult> Get([FromQuery] int page=1, [FromQuery] int limit=10, string status=null, [FromServices] ITaskService taskSvc) {
        var userId = int.Parse(User.FindFirst(ClaimTypes.NameIdentifier)?.Value);
        var role = User.FindFirst(ClaimTypes.Role)?.Value;
        var tasks = await taskSvc.GetTasks(userId, role, page, limit, status);
        return Ok(tasks);
    }
}
