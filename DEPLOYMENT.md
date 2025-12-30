# PixelScribe Deployment Guide (100% FREE)

> Complete guide to deploying PixelScribe backend (Go API) and frontend (React + Vite) for **$0/month** using best-in-class free tiers.

## 📋 Table of Contents

- [Architecture Overview](#architecture-overview)
- [Prerequisites](#prerequisites)
- [Quick Start](#quick-start)
- [Detailed Setup](#detailed-setup)
  - [Step 1: Database Setup (Neon.tech)](#step-1-database-setup-neontech)
  - [Step 2: Backend Deployment (Render)](#step-2-backend-deployment-render)
  - [Step 3: Frontend Deployment (Vercel)](#step-3-frontend-deployment-vercel)
  - [Step 4: CORS Configuration](#step-4-cors-configuration-backend)
- [Configuration Reference](#configuration-reference)
- [Cold Start Mitigation](#cold-start-mitigation-optional)
- [Monitoring & Debugging](#monitoring--logs)
- [Custom Domains](#custom-domain-optional)
- [Cost Analysis](#cost-comparison)
- [Troubleshooting](#troubleshooting)
- [Performance Optimization](#performance-optimization)

## 🏗️ Architecture Overview

```
┌─────────────────────────────────────────────────────────┐
│                    PixelScribe Stack                     │
├─────────────────────────────────────────────────────────┤
│                                                          │
│  ┌────────────────────┐         ┌──────────────────┐   │
│  │   Frontend (Vite)  │────────▶│  Backend (Go)    │   │
│  │   Vercel           │  HTTPS  │  Render          │   │
│  │   React + TS       │  API    │  Gin Framework   │   │
│  │   FREE, Always-on  │  Calls  │  FREE, Cold Start│   │
│  │   Global CDN       │         │  Docker Container│   │
│  └────────────────────┘         └─────────┬────────┘   │
│                                            │             │
│                                            │ PostgreSQL  │
│                                            ▼             │
│                                  ┌──────────────────┐   │
│                                  │  Database        │   │
│                                  │  Neon.tech       │   │
│                                  │  FREE, Serverless│   │
│                                  │  0.5GB Storage   │   │
│                                  └──────────────────┘   │
│                                                          │
│  Total Monthly Cost: $0                                 │
└─────────────────────────────────────────────────────────┘
```

## ✅ Prerequisites

Before deploying, ensure you have:

1. **GitHub Account** - To connect repositories
2. **Vercel Account** - [Sign up at vercel.com](https://vercel.com)
3. **Render Account** - [Sign up at render.com](https://render.com)
4. **Neon.tech Account** - [Sign up at neon.tech](https://neon.tech)
5. **OpenAI API Key** - [Get from platform.openai.com](https://platform.openai.com/api-keys)
6. **Repository Pushed to GitHub** - This project must be in a GitHub repo

### Local Tools (Optional, for CLI deployment)
```bash
# Install Vercel CLI (optional)
npm i -g vercel

# Install Render CLI (optional)
brew tap render-oss/render
brew install render
```

## 🚀 Quick Start

**TL;DR - Deploy in 15 minutes:**

1. **Database**: Neon.tech → Create project → Copy connection string
2. **Backend**: Render → New Web Service → Docker → Add env vars → Deploy
3. **Frontend**: Vercel → Import repo → Set `web` root → Add `VITE_API_URL` → Deploy

**Detailed instructions below** ⬇️

---

## 📖 Detailed Setup

### Step 1: Database Setup (Neon.tech)

**Why Neon?** Serverless PostgreSQL with permanent free tier (unlike Render's 90-day expiry).

1. **Create Account**
   - Go to [neon.tech](https://neon.tech)
   - Sign up with GitHub/Google/Email

2. **Create Project**
   - Click **"Create Project"**
   - Project Name: `pixelscribe`
   - Region: Choose closest to your users (US East recommended)
   - PostgreSQL Version: Latest (16+)
   - Click **"Create Project"**

3. **Get Connection String**
   - After project creation, you'll see a connection string
   - Format: `postgresql://user:pass@ep-xxx-xxx.us-east-2.aws.neon.tech/neondb?sslmode=require`
   - **Save this** - you'll need it for Render backend setup

4. **Verify Database**
   ```bash
   # Optional: Test connection locally
   psql "postgresql://user:pass@ep-xxx.neon.tech/neondb?sslmode=require"
   ```

**Free Tier Limits:**

- ✅ **0.5GB storage** (enough for 50K+ dictations)
- ✅ **100hrs compute/month** (plenty for MVP)
- ✅ **Auto-suspends when idle** (saves compute hours)
- ✅ **Permanent free tier** (no 90-day expiration)
- ✅ **Automatic backups** (7-day retention)

---

### Step 2: Backend Deployment (Render)

**Why Render?** Best free Docker hosting with automatic SSL and easy GitHub integration.

1. **Create Account & Connect GitHub**
   - Go to [render.com/dashboard](https://render.com/dashboard)
   - Sign up with GitHub
   - Authorize Render to access your repositories

2. **Create Web Service**
   - Click **"New +"** → **"Web Service"**
   - Select your `PixelScribe` repository
   - If not visible, click **"Configure account"** and grant access

3. **Configure Service**
   - **Name**: `pixelscribe-api` (or your preferred name)
   - **Region**: US West (Oregon) or closest to your users
   - **Branch**: `main` (or your default branch)
   - **Root Directory**: *(leave empty for root)*
   - **Environment**: `Docker`
   - **Dockerfile Path**: `./Dockerfile`
   - **Instance Type**: `Free`

4. **Add Environment Variables**

   Click **"Advanced"** → **"Add Environment Variable"** and add these:

   - **`DB_SOURCE`**: Your Neon connection string from Step 1
     - Format: `postgresql://user:pass@ep-xxx.neon.tech/neondb?sslmode=require`
   - **`SERVER_ADDRESS`**: `0.0.0.0:8080` (leave as-is)
   - **`TOKEN_SYMMETRIC_KEY`**: Generate a secure 32-character key (see below)
   - **`ACCESS_TOKEN_DURATION`**: `15m` (JWT token expiry)
   - **`OPENAI_API_KEY`**: Your OpenAI API key from platform.openai.com

   **Generate Secure TOKEN_SYMMETRIC_KEY:**
   ```bash
   # Run this locally to generate a secure 32-character key
   openssl rand -base64 32 | head -c 32

   # Example output: a4PTkR6Y6Ook8GkCvQElNO1FxqY9aZZf
   ```

5. **Deploy**
   - Click **"Create Web Service"**
   - Render will start building your Docker image
   - Initial build takes 5-10 minutes
   - Watch logs for any errors

6. **Verify Deployment**
   - Once deployed, you'll get a URL like: `https://pixelscribe-api.onrender.com`
   - Test health endpoint: `https://pixelscribe-api.onrender.com/health`
   - **Copy this URL** - you'll need it for frontend setup

**What Happens During Deployment:**

```bash
# Render automatically runs:
1. Builds Docker image from ./Dockerfile
2. Downloads dependencies (Go modules)
3. Compiles Go binary
4. Runs migrations (via start.sh entrypoint)
5. Starts the API server on port 8080
```

**Free Tier Limits:**

- ✅ **750 hours/month** (plenty for testing)
- ⚠️ **Sleeps after 15min inactivity** (cold start: 30-60s)
- ✅ **512MB RAM** (sufficient for Go API)
- ✅ **Automatic SSL/HTTPS**
- ✅ **Automatic deployments** on git push

---

### Step 3: Frontend Deployment (Vercel)

**Why Vercel?** Best-in-class static hosting with global CDN, instant deployments, and zero cold starts.

#### Option A: Using Vercel Dashboard (Recommended)

1. **Import Repository**
   - Go to [vercel.com/new](https://vercel.com/new)
   - Sign up/login with GitHub
   - Click **"Import Project"**
   - Select your **PixelScribe** repository

2. **Configure Project**
   - **Project Name**: `pixelscribe-web` (or your preference)
   - **Framework Preset**: Vite (auto-detected)
   - **Root Directory**: `web` ⚠️ **Important: Must be `web`**
   - **Build Command**: `npm run build` (auto-filled)
   - **Output Directory**: `dist` (auto-filled)
   - **Install Command**: `npm install` (auto-filled)

3. **Add Environment Variable**
   - Click **"Environment Variables"** section
   - Add variable:
     - **Name**: `VITE_API_URL`
     - **Value**: `https://pixelscribe-api.onrender.com` (your Render URL from Step 2)
     - **Environments**: Production, Preview, Development (select all)

4. **Deploy**
   - Click **"Deploy"**
   - Wait 1-2 minutes for build
   - You'll get a URL like: `https://pixelscribe-web.vercel.app`

5. **Test Deployment**
   - Visit your Vercel URL
   - Try registering a user
   - Test voice dictation feature
   - First API call may be slow (backend cold start)

#### Option B: Using Vercel CLI

```bash
# Install Vercel CLI globally
npm install -g vercel

# Navigate to web directory
cd web

# Deploy to production
vercel --prod

# Follow interactive prompts:
# ? Set up and deploy "~/PixelScribe/web"? Y
# ? Which scope? (Use arrow keys)
# ? Link to existing project? N
# ? What's your project's name? pixelscribe-web
# ? In which directory is your code located? ./
# ? Override settings? N

# Add environment variable
vercel env add VITE_API_URL production
# Paste: https://pixelscribe-api.onrender.com

# Redeploy with new env var
vercel --prod
```

**Free Tier Limits:**

- ✅ **100GB bandwidth/month** (plenty for MVP)
- ✅ **Unlimited deployments**
- ✅ **Global CDN** (sub-100ms worldwide)
- ✅ **Automatic HTTPS/SSL**
- ✅ **No cold starts** (always instant)
- ✅ **Automatic Git integration** (deploy on push)

---

### Step 4: CORS Configuration (Backend)

Your backend needs to allow requests from your Vercel domain. Check your CORS middleware:

**Location**: [internal/server/server.go](internal/server/server.go) (or wherever your CORS config is)

```go
// Make sure your CORS middleware includes your Vercel domain
config := cors.DefaultConfig()
config.AllowOrigins = []string{
    "http://localhost:5173",              // Local development
    "https://pixelscribe-web.vercel.app", // Your Vercel production URL
    "https://*.vercel.app",                // Vercel preview deployments
}
config.AllowCredentials = true
config.AllowHeaders = []string{"Origin", "Content-Type", "Authorization"}
router.Use(cors.New(config))
```

**If CORS is not configured**, you'll see errors like:

```text
Access to fetch at 'https://pixelscribe-api.onrender.com/api/...' from origin 'https://pixelscribe-web.vercel.app' has been blocked by CORS policy
```

**Fix**: Update your CORS config and redeploy backend (Render auto-deploys on git push).

---

### Step 5: Test Your Deployment

1. **Visit Frontend**
   - Open your Vercel URL: `https://pixelscribe-web.vercel.app`
   - Page should load instantly (no cold start)

2. **Test Registration**
   - Click "Sign Up" / "Register"
   - Create a new account
   - ⏱️ First API call takes 30-60s (backend waking up)
   - Subsequent calls are fast

3. **Test Core Features**
   - ✅ Login/Logout
   - ✅ Voice dictation (microphone access)
   - ✅ Text-to-speech playback
   - ✅ Save/retrieve dictations
   - ✅ OpenAI integration (if configured)

4. **Monitor Performance**
   - First load after 15min idle: Slow (cold start)
   - Active usage: Fast (<500ms API responses)
   - Frontend: Always instant (Vercel CDN)

---

## 🚀 Performance Optimization

### Cold Start Mitigation (Optional)

If 30-60s cold starts are unacceptable, you have options:

#### Option 1: Keep Backend Awake (FREE with cron service)

Use a free cron service to ping your backend every 10-14 minutes:

**Using cron-job.org (Free)**:

1. Go to [cron-job.org](https://cron-job.org)
2. Sign up for free account
3. Create new cron job:
   - **URL**: `https://pixelscribe-api.onrender.com/health`
   - **Schedule**: Every 10 minutes
   - **Method**: GET
4. Save and activate

**Using UptimeRobot (Free)**:

1. Go to [uptimerobot.com](https://uptimerobot.com)
2. Add new monitor:
   - **Monitor Type**: HTTP(s)
   - **URL**: `https://pixelscribe-api.onrender.com/health`
   - **Monitoring Interval**: 5 minutes (free tier)

**Trade-off**: Uses ~15 hours/month of your 750-hour free tier, but eliminates cold starts.

#### Option 2: Upgrade Backend Only ($7/month)

- Upgrade Render backend to paid plan: **$7/month**
- Keep frontend (Vercel) and database (Neon) **FREE**
- **Total cost**: $7/month with zero cold starts
- **Benefits**: Always-on, better performance, no hour limits

**How to upgrade**:

1. Render Dashboard → Your Service → Settings
2. Instance Type → Select "Starter" ($7/month)
3. Confirm upgrade

---

## 💰 Cost Comparison

| Setup                          | Monthly Cost | Cold Starts  | Storage | Best For                    |
|--------------------------------|--------------|--------------|---------|------------------------------|
| **100% Free (Recommended)**    | $0           | Yes (30-60s) | 0.5GB   | MVP, demos, low traffic     |
| **Free + Cron Keepalive**      | $0           | Minimal      | 0.5GB   | Testing, moderate traffic   |
| **Backend Paid**               | $7           | No           | 0.5GB   | Production, consistent load |
| **Backend + DB Paid**          | $26          | No           | 10GB    | Heavy usage, large dataset  |
| **AWS/GCP (comparison)**       | $20-50+      | No           | Custom  | Enterprise only             |

**Recommended Path**:

1. Start with **100% Free** ($0/month) for MVP and testing
2. Add **Free Cron Keepalive** ($0/month) if cold starts become annoying
3. Upgrade to **Backend Paid** ($7/month) when you get consistent users
4. Scale database as needed when storage exceeds 0.5GB

---

## 📊 Monitoring & Debugging

### Backend Logs (Render)

```bash
# View in dashboard
Render Dashboard → pixelscribe-api → Logs

# Common log checks:
# - Database connection success/failure
# - API endpoint hits
# - Migration status
# - Error traces
```

### Frontend Logs (Vercel)

```bash
# View in dashboard
Vercel Dashboard → pixelscribe-web → Deployments → View Function Logs

# Build logs
Vercel Dashboard → pixelscribe-web → Deployments → [Latest] → Building

# Runtime errors visible in browser console (F12)
```

### Database Monitoring (Neon)

```bash
# View in dashboard
Neon Dashboard → pixelscribe → Operations

# Check:
# - Active connections
# - Compute usage hours
# - Storage used
# - Query performance
```

### Health Check Endpoints

Test these endpoints to verify services:

```bash
# Backend health
curl https://pixelscribe-api.onrender.com/health

# Expected response: {"status": "ok"} or similar

# Frontend (just visit in browser)
https://pixelscribe-web.vercel.app
```

---

## 🌐 Custom Domain (Optional)

### Frontend Domain (Vercel)

1. **Add Domain to Vercel**
   - Vercel Dashboard → pixelscribe-web → Settings → Domains
   - Add domain: `pixelscribe.com`
   - Vercel provides DNS configuration

2. **Update DNS Records**
   - Add CNAME record: `pixelscribe.com` → `cname.vercel-dns.com`
   - Wait for propagation (5-60 minutes)
   - Vercel automatically provisions SSL certificate

### Backend Subdomain (Render)

1. **Add Custom Domain**
   - Render Dashboard → pixelscribe-api → Settings → Custom Domain
   - Add subdomain: `api.pixelscribe.com`

2. **Update DNS Records**
   - Add CNAME record: `api.pixelscribe.com` → `your-service.onrender.com`
   - SSL certificate auto-provisioned

3. **Update Frontend**
   - Update Vercel env var: `VITE_API_URL=https://api.pixelscribe.com`
   - Redeploy frontend

---

## 🔧 Configuration Reference

### Backend Environment Variables (Render)

| Variable               | Example Value                                                     | Required | Description                      |
|------------------------|-------------------------------------------------------------------|----------|----------------------------------|
| `DB_SOURCE`            | `postgresql://user:pass@ep-xxx.neon.tech/db?sslmode=require`     | Yes      | Neon PostgreSQL connection       |
| `SERVER_ADDRESS`       | `0.0.0.0:8080`                                                    | Yes      | Server bind address              |
| `TOKEN_SYMMETRIC_KEY`  | `a4PTkR6Y6Ook8GkCvQElNO1FxqY9aZZf`                                | Yes      | JWT signing key (32 chars)       |
| `ACCESS_TOKEN_DURATION`| `15m`                                                             | Yes      | JWT expiry time                  |
| `OPENAI_API_KEY`       | `sk-proj-xxxxx`                                                   | Optional | For AI features                  |

### Frontend Environment Variables (Vercel)

| Variable        | Example Value                          | Required | Description           |
|-----------------|----------------------------------------|----------|-----------------------|
| `VITE_API_URL`  | `https://pixelscribe-api.onrender.com` | Yes      | Backend API endpoint  |

---

## 🔍 Troubleshooting

### Problem: Backend Returns 502/503 Error

**Cause**: Backend is cold starting or crashed

**Solution**:

1. Check Render logs for errors
2. Wait 30-60 seconds and retry
3. Verify environment variables are set correctly
4. Check database connection in logs

### Problem: Frontend Shows CORS Error

**Error Message**:
```text
Access to fetch at 'https://api...' has been blocked by CORS policy
```

**Solution**:

1. Add Vercel domain to backend CORS config
2. Update [internal/server/server.go](internal/server/server.go)
3. Push changes to GitHub (Render auto-deploys)
4. Clear browser cache and retry

### Problem: Database Connection Failed

**Error Message**:
```text
connection to server failed: SSL required
```

**Solution**:

1. Verify `?sslmode=require` is in `DB_SOURCE`
2. Check Neon dashboard → Connection Details
3. Ensure connection string is complete
4. Test locally: `psql "your-connection-string"`

### Problem: Migrations Not Running

**Symptoms**: Tables don't exist, schema errors

**Solution**:

1. Check Render logs for migration output
2. Verify migration files exist in `/db/migration`
3. Check `start.sh` is executable
4. Manually run migrations:
   ```bash
   # In Render shell
   /app/migrate -path /app/db/migration -database "$DB_SOURCE" up
   ```

### Problem: OpenAI API Not Working

**Symptoms**: AI features fail, 401 errors

**Solution**:

1. Verify `OPENAI_API_KEY` is set in Render
2. Check API key is valid at platform.openai.com
3. Ensure billing is active on OpenAI account
4. Check Render logs for OpenAI errors

### Problem: Build Fails on Render

**Common Issues**:

1. **Go module errors**: Run `go mod tidy` locally and commit
2. **Missing dependencies**: Check `go.mod` is committed
3. **Docker build errors**: Test locally: `docker build -t test .`
4. **Migration tool download fails**: Check internet connectivity

### Problem: Vercel Build Fails

**Common Issues**:

1. **TypeScript errors**: Run `npm run build` locally first
2. **Missing dependencies**: Run `npm install` and commit `package-lock.json`
3. **Wrong root directory**: Ensure `web` is set as root
4. **Env var missing**: Add `VITE_API_URL` in Vercel dashboard

---

## 📈 Scaling & Upgrading

### When to Upgrade

**Indicators you've outgrown free tier**:

- Backend exceeds 750 hours/month (consistent 24/7 usage)
- Database exceeds 0.5GB storage
- Cold starts affecting user experience
- Need better uptime guarantees
- Require support/SLA

### Upgrade Path

1. **Phase 1: Eliminate Cold Starts** ($7/month)
   - Upgrade Render backend to Starter plan
   - Keep everything else free
   - Suitable for: 100-1000 active users

2. **Phase 2: Scale Database** ($26/month)
   - Upgrade Neon to Launch plan ($19/month)
   - 10GB storage, better performance
   - Suitable for: 1000-10000 users

3. **Phase 3: Professional Tier** ($50+/month)
   - Render Pro plan ($25/month)
   - Neon Scale plan ($69/month)
   - Multiple environments (staging/prod)
   - Suitable for: 10K+ users

---

## ✅ Deployment Checklist

### Pre-Deployment

- [ ] Code pushed to GitHub
- [ ] OpenAI API key obtained
- [ ] Accounts created (Neon, Render, Vercel)

### Database (Neon)

- [ ] Project created
- [ ] Connection string copied
- [ ] Database accessible (test with `psql`)

### Backend (Render)

- [ ] Web service created
- [ ] Docker environment selected
- [ ] All 5 environment variables added
- [ ] Build completed successfully
- [ ] Migrations ran successfully
- [ ] Health endpoint responding
- [ ] Backend URL copied

### Frontend (Vercel)

- [ ] Project imported from GitHub
- [ ] Root directory set to `web`
- [ ] `VITE_API_URL` environment variable added
- [ ] Build completed successfully
- [ ] Frontend URL accessible
- [ ] Can load homepage

### Integration Testing

- [ ] User registration works
- [ ] Login works
- [ ] Voice dictation works
- [ ] Text-to-speech works
- [ ] Dictations save to database
- [ ] CORS configured correctly
- [ ] No console errors

### Optional Optimizations

- [ ] Cron job set up (if avoiding cold starts)
- [ ] Custom domain configured (if applicable)
- [ ] Monitoring/alerts set up
- [ ] Error tracking configured

---

## 🎉 Success!

Your PixelScribe app is now deployed and accessible worldwide!

**What you've accomplished**:

- ✅ Professional-grade deployment on free tier
- ✅ Scalable architecture (frontend/backend/database separated)
- ✅ Automatic deployments (push to GitHub = auto-deploy)
- ✅ Global CDN for frontend (Vercel)
- ✅ Secure HTTPS for all services
- ✅ Database with automatic backups (Neon)

**Next steps**:

1. Share your frontend URL with users
2. Monitor usage in dashboards
3. Collect feedback and iterate
4. Upgrade when you hit free tier limits

**Need help?** Open an issue in the GitHub repository or check the main [README.md](../README.md).

---

**Made with ❤️ for PixelScribe** | [GitHub](https://github.com/yourusername/PixelScribe) | [Documentation](../README.md)
