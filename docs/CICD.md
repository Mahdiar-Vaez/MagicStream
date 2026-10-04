# راهنمای جامع خط‌به‌خط پایپ‌لاین CI/CD گیت‌هاب اکشن 🚀

این مستند توضیح کامل و خط‌به‌خط فایل کانفیگ `.github/workflows/deploy.yml` است که فرآیند خودکارسازی ساخت ایمیج‌ها، ارسال به داکر هاب و استقرار روی سرور را انجام می‌دهد.

---

## لیست سکرت‌های استفاده‌شده (GitHub Secrets)

این متغیرها باید در مسیر `Settings > Secrets and variables > Actions` در ریپازیتوری گیت‌هاب شما تنظیم شده باشند:

| نام Secret | توضیح |
|-------------|-------|
| `DOCKER_USERNAME` | نام کاربری حساب شما در Docker Hub |
| `DOCKER_PAT` | توکن دسترسی شخصی (Personal Access Token) داکر هاب برای لاگین ایمن |
| `SERVER_IP` | آدرس IP سرور مقصد (مثلاً `156.241.0.153`) |
| `SERVER_USER` | کاربر SSH سرور (مثلاً `root` یا `ubuntu`) |
| `SERVER_SSH_PORT` | پورت اتصال SSH (معمولاً `22`) |
| `SERVER_SSH_KEY` | کلید خصوصی SSH برای اتصال امن بدون پسورد |
| `SERVER_DIR` | مسیر قرارگیری پروژه روی سرور (مثلاً `/home/ubuntu/MagicStream`) |

---

## توضیح خط‌به‌خط کانفیگ `deploy.yml`

```yaml
1: name: CI/CD Pipeline - MagicStream
```
- **توضیح:** نام نمایشی پایپ‌لاین در تب **Actions** در گیت‌هاب.

```yaml
3: on:
4:   push:
5:     branches:
6:       - main
7:   workflow_dispatch:
```
- **خط ۳ تا ۶:** این بخش Trigger (ماشه) پایپ‌لاین را تعریف می‌کند؛ یعنی هر زمان که کدی روی شاخه اصلی (`main`) پوش (`push`) یا ادغام (`merge`) شود، پایپ‌لاین خودکار اجرا می‌شود.
- **خط ۷ (`workflow_dispatch`):** به شما اجازه می‌دهد در پنل گیت‌هاب در صورت تمایل دکمه **Run workflow** را بزنید و فرآیند را به صورت دستی اجرا کنید.

```yaml
9: jobs:
10:   build-and-push:
11:     name: Build & Push Docker Images
12:     runs-on: ubuntu-latest
```
- **خط ۱۰ و ۱۱:** تعریف اولین جاب (Job) تحت عنوان `build-and-push` برای ساخت و ارسال ایمیج‌ها به داکر هاب.
- **خط ۱۲:** محیط اجرایی جاب را تعیین می‌کند که یک ماشین مجازی تمیز و به‌روز اوبونتو (`ubuntu-latest`) ارائه‌شده توسط گیت‌هاب است.

```yaml
13:     steps:
14:       - name: Checkout Code
15:         uses: actions/checkout@v4
```
- **خط ۱۴ و ۱۵:** کدهای پروژه شما را از مخزن گیت‌هاب دانلود کرده و در فضای ماشین مجازی قرار می‌دهد تا داکر فایل‌ها و سورس‌کدها در دسترس باشند.

```yaml
17:       - name: Set up Docker Buildx
18:         uses: docker/setup-buildx-action@v3
```
- **خط ۱۷ و ۱۸:** ابزار پیشرفته `Docker Buildx` را فعال می‌کند. این ابزار امکاناتی نظیر کش کردن لایه‌ها (Caching) برای افزایش چشمگیر سرعت بیلد را فراهم می‌سازد.

```yaml
20:       - name: Log in to Docker Hub
21:         uses: docker/login-action@v3
22:         with:
23:           username: ${{ secrets.DOCKER_USERNAME }}
24:           password: ${{ secrets.DOCKER_PAT }}
```
- **خط ۲۰ تا ۲۴:** با استفاده از اکشن رسمی داکر و سکرت‌های ذخیره‌شده (`DOCKER_USERNAME` و `DOCKER_PAT`) وارد حساب کاربری شما در داکر هاب می‌شود تا مجوز پوش کردن ایمیج‌ها صادر شود.

```yaml
26:       - name: Build & Push API Image
27:         uses: docker/build-push-action@v6
28:         with:
29:           context: ./Server/MagicStreamServer
30:           file: ./Server/MagicStreamServer/Dockerfile.prod
31:           push: true
32:           tags: |
33:             ${{ secrets.DOCKER_USERNAME }}/magic-stream-api:latest
34:             ${{ secrets.DOCKER_USERNAME }}/magic-stream-api:${{ github.sha }}
35:           cache-from: type=gha
36:           cache-to: type=gha,mode=max
```
- **خط ۲۹:** کانتکست بیلد بک‌اند را به مسیر `./Server/MagicStreamServer` تنظیم می‌کند.
- **خط ۳۰:** صراحتاً داکرفایل پروداکشن بک‌اند (`Dockerfile.prod`) را مشخص می‌کند.
- **خط ۳۱ (`push: true`):** مشخص می‌کند که پس از پایان بیلد، ایمیج بلافاصله به داکر هاب پوش شود.
- **خط ۳۲ تا ۳۴:** دو تگ به ایمیج می‌زند:
  1. `latest`: نشان‌دهنده آخرین نسخه پایدار است که سرور آن را پول خواهد کرد.
  2. `${{ github.sha }}`: هش کامیت اختصاصی گیت (برای امکان بازگشت یا رول‌بک به یک نسخه خاص در آینده).
- **خط ۳۵ و ۳۶:** کش گیت‌هاب اکشن (`type=gha`) را فعال می‌کند تا پکیج‌های Go و لایه‌های تغییرنیافته مجدداً دانلود نشوند و زمان بیلد به زیر ۳۰ ثانیه برسد.

```yaml
38:       - name: Build & Push Client Image
39:         uses: docker/build-push-action@v6
40:         with:
41:           context: ./Client/magic-stream-client
42:           file: ./Client/magic-stream-client/Dockerfile.prod
43:           build-args: |
44:             VITE_API_BASE_URL=/api
45:           push: true
46:           tags: |
47:             ${{ secrets.DOCKER_USERNAME }}/magic-stream-client:latest
48:             ${{ secrets.DOCKER_USERNAME }}/magic-stream-client:${{ github.sha }}
49:           cache-from: type=gha
50:           cache-to: type=gha,mode=max
```
- **خط ۴۱ و ۴۲:** مسیر و داکرفایل پروداکشن فرانت‌اند React را مشخص می‌کند.
- **خط ۴۳ و ۴۴ (`build-args`):** مقدار حیاتی `VITE_API_BASE_URL=/api` را به فرآیند بیلد تزریق می‌کند تا فرانت‌اند به صورت کاملاً مستقل از دامنه/پورت، درخواست‌ها را به پروکسی داخلی `/api` بفرستد.
- **خط ۴۵ تا ۵۰:** ایمیج کلاینت را با تگ‌های `latest` و شناسه کامیت به داکر هاب ارسال کرده و لایه‌ها را کش می‌کند.

---

```yaml
52:   deploy:
53:     name: Deploy to Server via SSH
54:     runs-on: ubuntu-latest
55:     needs: build-and-push
```
- **خط ۵۲ تا ۵۴:** تعریف جاب دوم به نام `deploy` برای اجرای عملیات روی سرور لینوکس.
- **خط ۵۵ (`needs: build-and-push`):** وابستگی مهم؛ استقرار فقط و فقط در صورتی آغاز می‌شود که مرحله قبل (بیلد و پوش ایمیج‌ها) با ۱۰۰٪ موفقیت به اتمام رسیده باشد.

```yaml
56:     steps:
57:       - name: Execute Remote SSH Commands
58:         uses: appleboy/ssh-action@v1.0.3
59:         with:
60:           host: ${{ secrets.SERVER_IP }}
61:           username: ${{ secrets.SERVER_USER }}
62:           key: ${{ secrets.SERVER_SSH_KEY }}
63:           port: ${{ secrets.SERVER_SSH_PORT }}
```
- **خط ۵۷ تا ۶۳:** از اکشن امن و استاندارد `appleboy/ssh-action` استفاده می‌کند و با اطلاعات `SERVER_IP`, `SERVER_USER`, `SERVER_SSH_PORT` و کلید خصوصی `SERVER_SSH_KEY` یک ارتباط ایمن SSH با سرور برقرار می‌کند.

```yaml
64:           script: |
65:             set -e
```
- **خط ۶۵ (`set -e`):** یک دستور شل برای امنیت بیشتر؛ اگر هر کدام از دستورات بعدی با خطا مواجه شود، اجرای اسکریپت فوراً متوقف می‌شود و در گیت‌هاب وضعیت Fail اعلام می‌گردد.

```yaml
66:             echo "🚀 Connecting to server directory..."
67:             cd ${{ secrets.SERVER_DIR }}
```
- **خط ۶۶ و ۶۷:** ورود به پوشه پروژه روی سرور (مثلاً `/home/ubuntu/MagicStream`).

```yaml
69:             echo "📥 Fetching latest configs from Git..."
70:             git pull origin main
```
- **خط ۶۹ و ۷۰:** آخرین تغییرات احتمالی فایل‌های کانفیگ (`docker-compose.prod.yaml`، اسکریپت‌های سید و ...) را مستقیماً از گیت‌هاب دریافت می‌کند (فایل محرمانه `.env.prod` دست‌نخورده باقی می‌ماند).

```yaml
72:             echo "🐳 Pulling latest Docker images..."
73:             docker compose --env-file .env.prod -f docker-compose.prod.yaml pull
```
- **خط ۷۲ و ۷۳:** ایمیج‌های تازه‌ساخته‌شده نسخه `latest` هر دو سرویس API و Client را از داکر هاب دریافت می‌کند (سرور اصلاً درگیر فرایند سنگین بیلد نمی‌شود).

```yaml
75:             echo "🔄 Recreating containers..."
76:             docker compose --env-file .env.prod -f docker-compose.prod.yaml up -d --remove-orphans
```
- **خط ۷۵ و ۷۶:** کانتینرهای قدیمی را متوقف کرده و با ایمیج‌های جدید راه‌اندازی می‌کند. فلگ `--remove-orphans` کانتینرهای بدون استفاده قبلی را نیز پاکسازی می‌کند. دیتابیس مونگو بدون قطعی به کار خود ادامه می‌دهد.

```yaml
78:             echo "🧹 Cleaning up old unused images..."
79:             docker image prune -f
```
- **خط ۷۸ و ۷۹:** ایمیج‌های قدیمی و بدون استفاده قبلی داکر را از روی دیسک سرور حذف می‌کند تا حافظه سرور پر نشود.

```yaml
81:             echo "✅ Checking running containers status..."
82:             docker compose --env-file .env.prod -f docker-compose.prod.yaml ps
```
- **خط ۸۱ و ۸۲:** وضعیت سلامت و اجرای کانتینرها را در کنسول خروجی گیت‌هاب اکشن چاپ می‌کند تا بتوانید بلافاصله وضعیت سبز بودن سرور را بررسی کنید.
