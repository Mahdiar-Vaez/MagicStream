# MagicStream 🎬✨

پلتفرم استریم فیلم با سیستم پیشنهاد هوش مصنوعی، ساخته‌شده با React، Go (gin-gonic) و MongoDB.

---

## فهرست مطالب

- [درباره پروژه](#درباره-پروژه)
- [تکنولوژی‌ها](#تکنولوژیها)
- [ساختار پروژه](#ساختار-پروژه)
- [پیش‌نیازها](#پیشنیازها)
- [محیط توسعه (Dev)](#محیط-توسعه-dev)
- [استقرار پروداکشن (Prod)](#استقرار-پروداکشن-prod)
  - [نکته مهم: VITE_API_BASE_URL](#نکته-مهم-vite_api_base_url)
  - [راهنمای گام‌به‌گام برای سرور](#راهنمای-گامبهگام-برای-سرور)
  - [راهنمای .env.prod](#راهنمای-envprod)
- [بهترین روش دیپلوی](#بهترین-روش-دیپلوی)
- [مدیریت تگ‌های Docker Hub](#مدیریت-تگهای-docker-hub)
- [نکات امنیتی](#نکات-امنیتی)
- [رفع اشکال](#رفع-اشکال)
- [لینک ویدیوی آموزشی](#لینک-ویدیوی-آموزشی)

---

## درباره پروژه

این پروژه یک شبیه‌سازی کامل از یک پلتفرم مدرن استریم فیلم است که نشان می‌دهد چطور می‌توان تکنولوژی‌های مختلف را برای ساختن یک اپلیکیشن مقیاس‌پذیر و هوشمند کنار هم گذاشت.

- **فرانت‌اند React** — تجربه کاربری مدرن با React-Player
- **بک‌اند Go (gin-gonic)** — API با عملکرد بالا
- **موتور پیشنهاد هوش مصنوعی** — با LangChainGo و OpenAI
- **پایگاه داده MongoDB** — ذخیره‌سازی متادیتای فیلم‌ها و کاربران

---

## تکنولوژی‌ها

| لایه | تکنولوژی |
|------|-----------|
| Frontend / Client | JavaScript / React / Vite |
| Backend / Server | Go / gin-gonic |
| Storage / Database | MongoDB 7 |
| Container | Docker / Docker Compose |
| Web Server (Prod) | Nginx (Alpine) |

---

## ساختار پروژه

```
MagicStream/
├── Client/
│   └── magic-stream-client/     # React + Vite SPA
│       ├── Dockerfile.dev
│       ├── Dockerfile.prod       # Multi-stage: Node build + Nginx serve
│       └── nginx.prod.conf       # Nginx با SPA routing و cache headers
├── Server/
│   └── MagicStreamServer/        # Go API (gin-gonic)
│       ├── Dockerfile.dev        # Hot-reload با Air
│       └── Dockerfile.prod       # Multi-stage: Alpine build + minimal runtime
├── seed-data/                    # داده‌های اولیه MongoDB (JSON)
├── seed.js                       # اسکریپت seed (فقط یک‌بار اجرا می‌شود)
├── docker-compose.dev.yaml       # محیط توسعه با Watch
├── docker-compose.prod.yaml      # محیط پروداکشن از Docker Hub
├── .env                          # متغیرهای توسعه (در Git ignore است)
├── .env.prod.example             # نمونه .env.prod — این را کپی کن
└── .gitignore
```

---

## پیش‌نیازها

| ابزار | نسخه توصیه‌شده |
|-------|----------------|
| Docker | 26+ |
| Docker Compose | 2.22+ (برای Watch در dev) |
| Git | هر نسخه |

> **توجه:** برای پروداکشن فقط Docker و Docker Compose لازم است — نه Go، نه Node.

---

## محیط توسعه (Dev)

### راه‌اندازی سریع

```bash
# کلون کن
git clone https://github.com/Mahdiar-Vaez/MagicStream.git
cd MagicStream

# فایل env توسعه وجود دارد — نیازی به تغییر ندارد (مقادیر پیش‌فرض کافی‌اند)
# cat .env

# اجرا با Hot-Reload (Docker Compose Watch)
docker compose --env-file .env -f docker-compose.dev.yaml up --build --watch
```

| سرویس | آدرس |
|--------|-------|
| Client (React/Vite) | http://localhost:5173 |
| API (Go/gin) | http://localhost:8080 |

### نحوه کار Watch در Dev

- تغییرات **React** → sync به کانتینر Vite → HMR در مرورگر
- تغییرات **Go** → sync به کانتینر API → Air برنامه را ری‌بیلد می‌کند
- تغییر **Dockerfile یا go.mod/package.json** → کانتینر ری‌بیلد می‌شود

---

## استقرار پروداکشن (Prod)

### نکته مهم: VITE_API_BASE_URL

> ⚠️ **این مهم‌ترین نکته برای prod است — بدون توجه به این، اپ کار نمی‌کند.**

`VITE_API_BASE_URL` یک **Build-time argument** است، نه یک Runtime env var.  
Vite این مقدار را **داخل JavaScript کامپایل‌شده** می‌گذارد. بعد از build ایمیج، دیگر قابل تغییر نیست.

**این یعنی:**
- اگر `20071386/magic-stream-client:1.0.0` را با آدرس `http://localhost:8080` build کرده‌ای، کلاینت همیشه به localhost وصل می‌شود — نه به سرور واقعی.
- برای ست کردن آدرس درست، باید ایمیج را با `--build-arg` ری‌بیلد و push کنی.

**چک کردن وضعیت فعلی ایمیج:**
```bash
# بررسی می‌کنیم VITE_API_BASE_URL چه مقداری دارد
docker run --rm 20071386/magic-stream-client:1.0.0 \
  sh -c "grep -r 'localhost:8080\|http://api' /usr/share/nginx/html --include='*.js' -l"
```

اگر نام فایل‌هایی چاپ شد، URL داخل آن‌ها `localhost:8080` است و باید ری‌بیلد کنی:

```bash
# روی ماشین build (نه سرور) — دامنه API خود را بگذار
cd Client/magic-stream-client

docker build \
  -f Dockerfile.prod \
  --build-arg VITE_API_BASE_URL=http://YOUR_SERVER_IP:8080 \
  -t 20071386/magic-stream-client:1.0.0 \
  .

docker push 20071386/magic-stream-client:1.0.0
```

> اگر دامنه داری (مثلاً api.magicstream.com)، IP را با دامنه جایگزین کن.

---

### راهنمای گام‌به‌گام برای سرور

این دستورات را روی سرور لینوکسی اجرا کن:

#### ۱. اتصال به سرور و نصب Docker

```bash
# اتصال SSH
ssh user@YOUR_SERVER_IP

# نصب Docker (Ubuntu/Debian)
curl -fsSL https://get.docker.com | sh
sudo usermod -aG docker $USER
newgrp docker
```

#### ۲. کپی فایل‌های لازم به سرور

فقط این فایل‌ها روی سرور لازم است — کد سورس نه، ایمیج‌ها از Docker Hub pull می‌شوند:

```bash
# روی ماشین local — ایجاد پوشه روی سرور و کپی فایل‌ها
ssh user@YOUR_SERVER_IP "mkdir -p ~/magicstream"

scp docker-compose.prod.yaml user@YOUR_SERVER_IP:~/magicstream/
scp .env.prod.example user@YOUR_SERVER_IP:~/magicstream/
scp seed.js user@YOUR_SERVER_IP:~/magicstream/
scp -r seed-data user@YOUR_SERVER_IP:~/magicstream/
```

#### ۳. ساختن فایل .env.prod روی سرور

```bash
# روی سرور
cd ~/magicstream
cp .env.prod.example .env.prod
nano .env.prod   # یا vim .env.prod
```

#### ۴. اجرای پروداکشن

```bash
# روی سرور — از پوشه magicstream
docker compose --env-file .env.prod -f docker-compose.prod.yaml pull
docker compose --env-file .env.prod -f docker-compose.prod.yaml up -d --remove-orphans
```

#### ۵. بررسی وضعیت

```bash
# وضعیت کانتینرها
docker compose --env-file .env.prod -f docker-compose.prod.yaml ps

# لاگ‌ها
docker compose --env-file .env.prod -f docker-compose.prod.yaml logs -f

# health check
docker inspect --format='{{.State.Health.Status}}' magicstream-prod-api-1
```

---

### راهنمای .env.prod

فایل `.env.prod.example` را کپی کن و مقادیر زیر را پر کن:

```ini
# پورت‌های expose‌شده روی سرور (پورت host)
PROD_CLIENT_PORT=80          # کلاینت روی این پورت در دسترس است
PROD_API_PORT=8080           # API روی این پورت در دسترس است

# Origin مجاز برای CORS
# اگر کلاینت روی http://1.2.3.4 (بدون پورت) است:
PROD_ALLOWED_ORIGINS=http://YOUR_SERVER_IP
# اگر پورت غیر ۸۰ است:
# PROD_ALLOWED_ORIGINS=http://YOUR_SERVER_IP:3000
# اگر دامنه داری:
# PROD_ALLOWED_ORIGINS=https://magicstream.example.com

# کلیدهای JWT — باید طولانی و تصادفی باشند (حداقل ۳۲ کاراکتر)
PROD_SECRET_KEY=replace-with-a-long-random-secret
PROD_SECRET_REFRESH_KEY=replace-with-a-different-long-random-secret

# OpenAI (اختیاری — اگر نداری خالی بگذار)
OPENAI_API_KEY=
BASE_PROMPT_TEMPLATE=
RECOMMENDED_MOVIE_LIMIT=5
```

> **تولید Secret Key امن:**
> ```bash
> # لینوکس/Mac
> openssl rand -hex 32
> # یا
> python3 -c "import secrets; print(secrets.token_hex(32))"
> ```

#### جدول خطاهای رایج .env.prod

| مشکل | علت | راه‌حل |
|-------|------|---------|
| API کار نمی‌کند | `PROD_SECRET_KEY` خالی یا کوتاه است | مقدار قوی وارد کن |
| CORS error در مرورگر | `PROD_ALLOWED_ORIGINS` اشتباه است | IP یا دامنه دقیق فرانت‌اند را وارد کن |
| کلاینت به API وصل نمی‌شود | `VITE_API_BASE_URL` اشتباه در ایمیج build شده | ایمیج کلاینت را ری‌بیلد کن (بالا توضیح داده شد) |
| seed اجرا نمی‌شود | Volume از قبل وجود دارد | Volume را حذف کن: `docker volume rm magicstream-prod-mongodb-data` |

---

## بهترین روش دیپلوی

### گزینه ۱ — ساده‌ترین (بدون دامنه) ✅ توصیه برای شروع

اجرای مستقیم با Docker Compose روی IP سرور:
- کلاینت: `http://SERVER_IP:80`
- API: `http://SERVER_IP:8080`

```bash
docker compose --env-file .env.prod -f docker-compose.prod.yaml up -d
```

---

### گزینه ۲ — با Nginx Reverse Proxy (توصیه برای پروداکشن واقعی) ⭐

نصب Nginx روی سرور و هدایت ترافیک:

```
Browser → Nginx (port 80/443) → Client container (port 80)
                              → API container (port 8080)  [path /api]
```

**مزایا:**
- HTTPS/TLS با Certbot رایگان
- مخفی کردن پورت API (8080) از اینترنت
- Rate limiting و security headers

```nginx
# /etc/nginx/sites-available/magicstream
server {
    listen 80;
    server_name magicstream.example.com;

    # Client
    location / {
        proxy_pass http://localhost:80;
        proxy_set_header Host $host;
        proxy_set_header X-Real-IP $remote_addr;
    }

    # API
    location /api/ {
        proxy_pass http://localhost:8080/;
        proxy_set_header Host $host;
        proxy_set_header X-Real-IP $remote_addr;
    }
}
```

بعد از تنظیم Nginx، HTTPS رایگان بگیر:
```bash
sudo apt install certbot python3-certbot-nginx
sudo certbot --nginx -d magicstream.example.com
```

---

### گزینه ۳ — با Traefik (پیشرفته)

اگر چند سرویس داری، Traefik به‌عنوان reverse proxy هوشمند داخل Docker Compose جالب است.

---

## مدیریت تگ‌های Docker Hub

**تگ‌های فعلی روی Docker Hub:**
- `20071386/magic-stream-api:1.0.0`
- `20071386/magic-stream-client:1.0.0`

### انتشار نسخه جدید

```bash
# ۱. بیلد و push ایمیج API
cd Server/MagicStreamServer
docker build -f Dockerfile.prod -t 20071386/magic-stream-api:1.1.0 .
docker push 20071386/magic-stream-api:1.1.0

# ۲. بیلد و push ایمیج Client (با آدرس API)
cd ../../Client/magic-stream-client
docker build -f Dockerfile.prod \
  --build-arg VITE_API_BASE_URL=http://YOUR_SERVER_IP:8080 \
  -t 20071386/magic-stream-client:1.1.0 .
docker push 20071386/magic-stream-client:1.1.0

# ۳. تگ در docker-compose.prod.yaml را به‌روزرسانی کن
# api: image: 20071386/magic-stream-api:1.1.0
# client: image: 20071386/magic-stream-client:1.1.0

# ۴. روی سرور آپدیت کن
docker compose --env-file .env.prod -f docker-compose.prod.yaml pull
docker compose --env-file .env.prod -f docker-compose.prod.yaml up -d
```

---

## نکات امنیتی

| ✅ انجام شده | ❌ انجام نشده |
|-------------|--------------|
| `.env` و `.env.prod` در `.gitignore` | — |
| `.env` در `.dockerignore` هر سرویس | — |
| Secret‌های prod فقط Runtime، نه داخل ایمیج | — |
| API با user غیر‌root اجرا می‌شود | — |
| GIN_MODE=release در prod | — |
| HTTPS (نیاز به Nginx + Certbot) | ← اگر دامنه داری |
| Firewall روی سرور (UFW) | ← توصیه می‌شود |

### فعال کردن Firewall (Ubuntu)

```bash
sudo ufw allow ssh
sudo ufw allow 80/tcp
sudo ufw allow 443/tcp
# پورت 8080 را فقط اگر مستقیم expose می‌کنی باز کن
sudo ufw allow 8080/tcp
sudo ufw enable
```

---

## رفع اشکال

### بررسی وضعیت کلی

```bash
docker compose --env-file .env.prod -f docker-compose.prod.yaml ps
docker compose --env-file .env.prod -f docker-compose.prod.yaml logs --tail=50
```

### کانتینر restart می‌شود (Restart Loop)

```bash
docker logs magicstream-prod-api-1 --tail=100
```

### بررسی config نهایی (بدون افشای secret)

```bash
# این خروجی را share نکن — شامل secret‌هاست
docker compose --env-file .env.prod -f docker-compose.prod.yaml config
```

### حذف کامل و شروع مجدد

```bash
# ⚠️ این Volume MongoDB را هم حذف می‌کند — داده‌ها از دست می‌روند
docker compose --env-file .env.prod -f docker-compose.prod.yaml down -v
docker compose --env-file .env.prod -f docker-compose.prod.yaml up -d
```

### تست health endpoint

```bash
curl http://localhost:8080/hello
# باید 200 برگرداند
```

---

## نکات مهم درباره فایل‌های محیطی

### .env (توسعه)
- **در Git نیست** (`.gitignore`)
- مقادیر پیش‌فرض dev دارد
- Compose این فایل را برای جایگزینی `${VARIABLE}` در YAML می‌خواند

### .env.prod (پروداکشن)  
- **هرگز commit نکن**
- **هرگز در Docker Hub push نکن**
- فقط روی سرور وجود دارد
- از `.env.prod.example` ساخته می‌شود

### VITE_API_BASE_URL — خاص است
- **Build-time** است، نه runtime
- Vite آن را داخل JS کامپایل می‌کند
- بعد از `docker build`، دیگر قابل تغییر نیست
- در `.env.prod` روی سرور **تأثیری ندارد**
- باید موقع `docker build` با `--build-arg` ست شود

---

## لینک ویدیوی آموزشی

- https://youtu.be/jBf7of9JTV8
