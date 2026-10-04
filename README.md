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
  - [معماری سرور](#معماری-سرور)
  - [راهنمای گام‌به‌گام برای سرور](#راهنمای-گامبهگام-برای-سرور)
  - [راهنمای .env.prod](#راهنمای-envprod)
- [نکته مهم: VITE_API_BASE_URL](#نکته-مهم-vite_api_base_url)
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

### معماری سرور

```
اینترنت
   │
   ├─── http://parva-ai.ir  (پورت 80)  ──► کانتینر client (Nginx)
   │
   └─── http://parva-ai.ir:8080        ──► کانتینر api (Go/gin)
                                              │
                                              └──► کانتینر db (MongoDB)
                                                   [فقط داخل Docker network]
```

> Docker images روی Docker Hub:
> - `20071386/magic-stream-api:1.0.0`
> - `20071386/magic-stream-client:1.0.0` ← build شده با `VITE_API_BASE_URL=http://parva-ai.ir:8080`

---

### راهنمای گام‌به‌گام برای سرور

#### ۱. اتصال به سرور و نصب Docker

```bash
# اتصال SSH
ssh root@156.241.0.153

# نصب Docker (Ubuntu/Debian)
curl -fsSL https://get.docker.com | sh

# اگر کاربر غیر root داری
sudo usermod -aG docker $USER
newgrp docker

# تست
docker --version
docker compose version
```

#### ۲. ایجاد پوشه پروژه روی سرور

```bash
mkdir -p ~/magicstream
cd ~/magicstream
```

#### ۳. کپی فایل‌های لازم به سرور

فقط این ۴ چیز روی سرور لازم است — کد سورس نه، ایمیج‌ها از Docker Hub pull می‌شوند:

**روش اول — با scp (از local):**
```bash
scp docker-compose.prod.yaml root@156.241.0.153:~/magicstream/
scp .env.prod.example root@156.241.0.153:~/magicstream/
scp seed.js root@156.241.0.153:~/magicstream/
scp -r seed-data root@156.241.0.153:~/magicstream/
```

**روش دوم — با git clone (روی سرور):**
```bash
# روی سرور
git clone https://github.com/Mahdiar-Vaez/MagicStream.git .
# فقط فایل‌های لازم کپی شدند — seed-data هم هست
```

#### ۴. ساختن فایل .env.prod روی سرور

```bash
cd ~/magicstream
cp .env.prod.example .env.prod
nano .env.prod
```

محتوای `.env.prod` را این‌طور پر کن (فقط secret‌ها را تغییر بده):

```ini
PROD_CLIENT_PORT=80
PROD_API_PORT=8080
PROD_ALLOWED_ORIGINS=http://parva-ai.ir
PROD_SECRET_KEY=اینجا-یک-کلید-بلند-تصادفی-بگذار
PROD_SECRET_REFRESH_KEY=اینجا-یک-کلید-بلند-تصادفی-دیگر-بگذار
OPENAI_API_KEY=
BASE_PROMPT_TEMPLATE=
RECOMMENDED_MOVIE_LIMIT=5
```

> **تولید secret key امن روی سرور:**
> ```bash
> openssl rand -hex 32
> # دو بار اجرا کن — یکی برای SECRET_KEY، یکی برای SECRET_REFRESH_KEY
> ```

#### ۵. اجرای پروداکشن

```bash
# pull ایمیج‌ها از Docker Hub
docker compose --env-file .env.prod -f docker-compose.prod.yaml pull

# اجرا در background
docker compose --env-file .env.prod -f docker-compose.prod.yaml up -d --remove-orphans
```

#### ۶. بررسی وضعیت

```bash
# وضعیت کانتینرها (همه باید Up باشند)
docker compose --env-file .env.prod -f docker-compose.prod.yaml ps

# لاگ‌ها
docker compose --env-file .env.prod -f docker-compose.prod.yaml logs -f

# تست API
curl http://localhost:8080/hello

# تست از بیرون
curl http://parva-ai.ir:8080/hello
curl http://parva-ai.ir
```

---

### راهنمای .env.prod

#### جدول متغیرها

| متغیر | مقدار برای این سرور | توضیح |
|--------|---------------------|-------|
| `PROD_CLIENT_PORT` | `80` | پورت host برای کلاینت |
| `PROD_API_PORT` | `8080` | پورت host برای API |
| `PROD_ALLOWED_ORIGINS` | `http://parva-ai.ir` | CORS origin — دقیق همین |
| `PROD_SECRET_KEY` | **باید تغییر کنی** | کلید JWT — حداقل ۳۲ کاراکتر |
| `PROD_SECRET_REFRESH_KEY` | **باید تغییر کنی** | کلید JWT refresh — متفاوت از بالا |
| `OPENAI_API_KEY` | اختیاری | اگر نداری خالی بگذار |
| `BASE_PROMPT_TEMPLATE` | اختیاری | قالب سفارشی prompt |
| `RECOMMENDED_MOVIE_LIMIT` | `5` | تعداد پیشنهاد AI |

#### ❌ این متغیر در .env.prod نیست و تأثیری ندارد

| متغیر | چرا اینجا نیست |
|--------|----------------|
| `VITE_API_BASE_URL` | Build-time است — داخل JS کامپایل شده. روی سرور قابل تغییر نیست |

#### جدول خطاهای رایج

| مشکل | علت احتمالی | راه‌حل |
|-------|-------------|---------|
| API کار نمی‌کند | `PROD_SECRET_KEY` کوتاه یا خالی | مقدار قوی وارد کن |
| CORS error در مرورگر | `PROD_ALLOWED_ORIGINS` اشتباه | دقیقاً `http://parva-ai.ir` باشد |
| کلاینت به API وصل نمی‌شود | ایمیج قدیمی با `localhost:8080` | ایمیج جدید از Docker Hub pull کن |
| صفحه سفید / 404 | Nginx SPA routing مشکل | لاگ کانتینر client را چک کن |
| seed اجرا نمی‌شود | Volume از قبل وجود دارد | این طبیعی است — seed فقط اول بار |

---

## نکته مهم: VITE_API_BASE_URL

> ⚠️ **این مهم‌ترین نکته برای prod است.**

`VITE_API_BASE_URL` یک **Build-time argument** است، نه Runtime env var.  
Vite این مقدار را **داخل JavaScript کامپایل‌شده** قرار می‌دهد. بعد از build ایمیج، دیگر قابل تغییر نیست.

**ایمیج فعلی روی Docker Hub:**
```
20071386/magic-stream-client:1.0.0
VITE_API_BASE_URL = http://parva-ai.ir:8080  ✅
```

### اگر سرور یا دامنه تغییر کند

```bash
# روی ماشین local — دامنه جدید را جایگزین کن
cd Client/magic-stream-client

docker build \
  -f Dockerfile.prod \
  --build-arg VITE_API_BASE_URL=http://DOMAIN_JADID:8080 \
  -t 20071386/magic-stream-client:1.0.0 .

docker push 20071386/magic-stream-client:1.0.0

# روی سرور — pull ایمیج جدید
docker compose --env-file .env.prod -f docker-compose.prod.yaml pull
docker compose --env-file .env.prod -f docker-compose.prod.yaml up -d
```

---

## مدیریت تگ‌های Docker Hub

**تگ‌های فعلی روی Docker Hub:**
- `20071386/magic-stream-api:1.0.0`
- `20071386/magic-stream-client:1.0.0` ← آدرس `http://parva-ai.ir:8080`

### انتشار نسخه جدید

```bash
# ۱. build و push API
cd Server/MagicStreamServer
docker build -f Dockerfile.prod -t 20071386/magic-stream-api:1.1.0 .
docker push 20071386/magic-stream-api:1.1.0

# ۲. build و push Client
cd ../../Client/magic-stream-client
docker build -f Dockerfile.prod \
  --build-arg VITE_API_BASE_URL=http://parva-ai.ir:8080 \
  -t 20071386/magic-stream-client:1.1.0 .
docker push 20071386/magic-stream-client:1.1.0

# ۳. تگ‌ها را در docker-compose.prod.yaml آپدیت کن
# ۴. روی سرور pull و restart کن
docker compose --env-file .env.prod -f docker-compose.prod.yaml pull
docker compose --env-file .env.prod -f docker-compose.prod.yaml up -d
```

---

## نکات امنیتی

| وضعیت | موضوع |
|--------|--------|
| ✅ | `.env` و `.env.prod` در `.gitignore` |
| ✅ | `.env` در `.dockerignore` هر سرویس |
| ✅ | Secret‌های prod فقط Runtime، نه داخل ایمیج |
| ✅ | API با user غیر‌root اجرا می‌شود |
| ✅ | `GIN_MODE=release` در prod |
| ⬜ | Firewall (UFW) روی سرور — توصیه می‌شود |
| ⬜ | HTTPS با Nginx + Certbot — اگر دامنه SSL داری |

### فعال کردن Firewall (Ubuntu)

```bash
sudo ufw allow ssh
sudo ufw allow 80/tcp
sudo ufw allow 8080/tcp
sudo ufw enable
sudo ufw status
```

### HTTPS رایگان (اختیاری — بعداً)

اگر خواستی HTTPS اضافه کنی:
```bash
sudo apt install nginx certbot python3-certbot-nginx -y
sudo certbot --nginx -d parva-ai.ir
```

---

## رفع اشکال

### بررسی وضعیت کلی

```bash
docker compose --env-file .env.prod -f docker-compose.prod.yaml ps
docker compose --env-file .env.prod -f docker-compose.prod.yaml logs --tail=50
```

### لاگ یک سرویس خاص

```bash
docker compose --env-file .env.prod -f docker-compose.prod.yaml logs api --tail=100
docker compose --env-file .env.prod -f docker-compose.prod.yaml logs client --tail=50
docker compose --env-file .env.prod -f docker-compose.prod.yaml logs db --tail=50
```

### restart یک سرویس

```bash
docker compose --env-file .env.prod -f docker-compose.prod.yaml restart api
```

### تست health endpoint

```bash
curl http://localhost:8080/hello
```

### حذف کامل و شروع مجدد (خطرناک!)

```bash
# ⚠️ این Volume MongoDB را هم حذف می‌کند — داده‌ها از دست می‌روند
docker compose --env-file .env.prod -f docker-compose.prod.yaml down -v
docker compose --env-file .env.prod -f docker-compose.prod.yaml up -d
```

---

## لینک ویدیوی آموزشی

- https://youtu.be/jBf7of9JTV8
