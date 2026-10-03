# MagicStream 🎬✨

Movie streaming platform with AI recommendation built with modern web technologies (React/Go/gin-gonic/MongoDB) 

---

## About  

This project is a full-stack simulation of a modern **Movie Streaming Platform**, designed to showcase how different technologies can be combined to deliver a scalable, AI-powered application.  

The system brings together a **React-based frontend** for an engaging user experience, a **Go-based backend** for high-performance API services that runs on the gin (gin-gonic) web framework, and an **AI-powered recommendation engine** to personalize movie suggestions using **LangChainGo** and **OpenAI**.  

It also demonstrates how **MongoDB** can serve as a reliable, scalable database solution for managing media metadata and user preferences.  

---

## Features

- Movie Streaming service simulated on the front end using React and React-Player
- Web API service written using GO and runs on the gin-gonic web framework 
- AI Recommendation service using LangChainGo, Go and OpenAI
- Scalable backend storage provided by MongoDB

---

## Tech Stack

| Frontend / Client | JavaScript / React |
| Backend / Server | Go / gin-gonic |
| Storage / Database | MongoDB |
 
---

## Link to Video Tutorial on How to Build the App
- https://youtu.be/jBf7of9JTV8

---

### Installation

1. Clone the repo  
   ```bash
   git clone https://github.com/GavinLonDigital/MagicStream.git
   cd MagicStream
   ```

## Docker

### فایل‌های محیطی و build داکر

نیازی به یک فایل `.env` سراسری برای Docker نیست. در این پروژه، Compose فایل `.env` ریشهٔ مخزن را می‌خواند تا مقدارهای `${VARIABLE}` را در فایل YAML جایگزین کند. این کار مستقل از فرستادن فایل به فرایند build ایمیج است.

هر سرویس build context جداگانه دارد: `Server/MagicStreamServer` برای API و `Client/magic-stream-client` برای کلاینت. فایل `.dockerignore` نسبت به build context همان سرویس اعمال می‌شود؛ فایل‌های matchشده را از ورودی `docker build` حذف می‌کند، اما جلوی خواندن `.env` ریشه توسط Compose را نمی‌گیرد. `.dockerignore`های API و کلاینت، `.env` را حذف می‌کنند تا اطلاعات محلی و محرمانه وارد ایمیج نشود.

تنظیمات API از بخش `environment` در Compose به‌صورت متغیر محیطی و هنگام اجرای کانتینر داده می‌شوند. Compose مقدار secretهای production را از `.env.prod` می‌خواند و به کانتینر API می‌دهد؛ این مقادیر build argument نیستند و نباید داخل ایمیج کپی شوند. کلاینت متفاوت است: Vite هنگام `npm run build` مقدار `VITE_API_BASE_URL` را داخل JavaScript نهایی قرار می‌دهد. پس آدرس عمومی API باید موقع build کردن image کلاینت تنظیم شود؛ `.env.prod` روی سرور نمی‌تواند image از قبل ساخته‌شده را تغییر دهد. این آدرس برای مرورگر قابل‌مشاهده است و secret محسوب نمی‌شود. هیچ API key یا رمزی را در متغیرهای `VITE_*` قرار نده.

برای توسعه از `.env` ریشه استفاده کن و فرمان Compose را از ریشهٔ مخزن اجرا کن. برای انتخاب فایل دیگری، مسیرش را با `--env-file` مشخص کن:

```bash
docker compose --env-file .env -f docker-compose.dev.yaml up --build --watch
```

برای production، فایل جداگانهٔ `.env.prod` روی سرور داشته باش و آن را به Git یا Docker Hub اضافه نکن. فایل را از روی نمونه بساز، دامنهٔ واقعی فرانت‌اند و secretهای نمونه را جایگزین کن و آن را صریح به Compose بده:

```bash
cp .env.prod.example .env.prod
docker compose --env-file .env.prod -f docker-compose.prod.yaml pull
docker compose --env-file .env.prod -f docker-compose.prod.yaml up -d --no-build --remove-orphans
```

در PowerShell، دستور کپی فایل این است: `Copy-Item .env.prod.example .env.prod`.

فایل‌های واقعی `.env` و `.env.prod` را وارد build context نکن و در Dockerfile کپی نکن. secretهای production را commit نکن. برای بررسی تنظیمات نهایی می‌توانی از `docker compose config` استفاده کنی؛ اگر خروجی آن secret دارد، آن را منتشر نکن.

محیط توسعه از Docker Compose Watch استفاده می‌کند. تغییرات React به کانتینر Vite همگام می‌شوند و HMR آن‌ها را در مرورگر اعمال می‌کند؛ تغییرات سورس Go به کانتینر API منتقل می‌شوند و Air برنامه را دوباره build و راه‌اندازی می‌کند. با تغییر Dockerfile یا فایل‌های وابستگی، Compose ایمیج مربوط را دوباره می‌سازد. این قابلیت به Docker Compose نسخهٔ 2.22 یا جدیدتر نیاز دارد:

```bash
docker compose -f docker-compose.dev.yaml up --build --watch
```

کلاینت در `http://localhost:5173` و API در `http://localhost:8080` در دسترس است.

برای production، فایل `.env.prod.example` را به `.env.prod` کپی کن، آدرس عمومی API، مبدأ مجاز فرانت‌اند و secretهای جداگانه و قوی را تنظیم کن، سپس اجرا کن:

```bash
docker compose --env-file .env.prod -f docker-compose.prod.yaml up --build -d
```

Compose در production imageهای `20071386/magic-stream-api:1.0.0` و `20071386/magic-stream-client:1.0.0` را از Docker Hub می‌گیرد و روی سرور build نمی‌کند. برای انتشار نسخهٔ جدید، imageها را با tag جدید push کن و همان tagها را در Compose به‌روزرسانی کن. اگر کلاینت قبلی را بدون build arg ساختی، احتمالاً API URL آن `http://localhost:8080` است و باید قبل از deploy دوباره build و push شود. این کار را روی سیستم build انجام بده؛ URL زیر را با دامنهٔ واقعی API جایگزین کن:

```powershell
docker build -f Dockerfile.prod --build-arg VITE_API_BASE_URL=https://api.example.com -t 20071386/magic-stream-client:1.0.0 .
docker push 20071386/magic-stream-client:1.0.0
```

کلاینت production به‌طور پیش‌فرض روی پورت 80 و API روی پورت 8080 ارائه می‌شوند. برای تغییر پورت‌های میزبان، `PROD_CLIENT_PORT` و `PROD_API_PORT` را در `.env.prod` تنظیم کن. محیط‌های dev و prod از volumeهای جداگانهٔ MongoDB استفاده می‌کنند. سرور علاوه بر Compose و `.env.prod`، باید `seed.js` و پوشهٔ `seed-data` را هم داشته باشد؛ seed فقط هنگام مقداردهی اولیهٔ volume خالی اجرا می‌شود.
