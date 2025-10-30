# Complete Building Guide - Shared Expenses App

This guide covers building and running both the backend and frontend applications.

## Prerequisites

### Backend Requirements
- **Go**: 1.21 or later
- **PostgreSQL**: 14 or later (running and accessible)
- **Environment Variables**: See backend/.env.default

### Frontend Requirements
- **Flutter SDK**: 3.9.2 or later
- **FVM**: For Flutter version management (recommended)
- **Android Studio** (for Android)
- **Xcode** (for iOS, macOS only)

---

## Part 1: Backend Setup

### 1.1 Install Dependencies

```bash
cd backend
go mod download
```

### 1.2 Configure Environment

```bash
# Copy the default environment file
cp .env.default .env

# Edit .env with your database credentials
nano .env
```

**Required settings:**
```env
DB_URL=postgres://username:password@localhost:5432/shared_expenses
API_PORT=8080
DB_MIGRATIONS_DIR=db/migrations
```

### 1.3 Setup Database

```bash
# Create PostgreSQL database
psql -U postgres
CREATE DATABASE shared_expenses;
\q

# Migrations run automatically on app start
```

### 1.4 Run Backend Server

```bash
# From backend directory
go run main.go

# Or build and run
go build -o shared-expenses-server
./shared-expenses-server
```

**Expected output:**
```
Server running on port 8080
```

**Test backend:**
```bash
curl http://localhost:8080/health
# Should return: ok
```

---

## Part 2: Frontend Setup

### 2.1 Install Flutter via FVM

```bash
# Install FVM globally
dart pub global activate fvm

# Navigate to frontend directory
cd frontend

# Install Flutter version
fvm install

# Use the configured Flutter version
fvm use stable
```

### 2.2 Install Dependencies

```bash
# Get all Flutter packages
fvm flutter pub get
```

### 2.3 Configure API Endpoint

Edit `lib/config/api_config.dart`:

**For Android Emulator (default):**
```dart
static const String baseUrl = 'http://10.0.2.2:8080';
```

**For iOS Simulator:**
```dart
static const String baseUrl = 'http://localhost:8080';
```

**For Physical Device:**
```dart
// Replace with your computer's local IP
static const String baseUrl = 'http://192.168.1.xxx:8080';
```

**Find your IP:**
```bash
# Linux/macOS
ifconfig | grep inet

# Windows
ipconfig
```

### 2.4 Run Flutter App

#### Option A: Android Emulator

```bash
# List available emulators
fvm flutter emulators

# Launch an emulator
fvm flutter emulators --launch <emulator-name>

# Run the app
fvm flutter run
```

#### Option B: iOS Simulator (macOS only)

```bash
# Open simulator
open -a Simulator

# Run the app
fvm flutter run
```

#### Option C: Physical Device

```bash
# Connect device via USB
# Enable USB debugging (Android) or Developer Mode (iOS)

# Check device is connected
fvm flutter devices

# Run on device
fvm flutter run
```

---

## Building for Production

### Backend Production Build

```bash
cd backend

# Build for Linux
GOOS=linux GOARCH=amd64 go build -o shared-expenses-server

# Build for Windows
GOOS=windows GOARCH=amd64 go build -o shared-expenses-server.exe

# Build for macOS
GOOS=darwin GOARCH=amd64 go build -o shared-expenses-server
```

### Frontend Production Builds

#### Android APK

```bash
cd frontend

# Build release APK
fvm flutter build apk --release

# Output: build/app/outputs/flutter-apk/app-release.apk

# Install on device
fvm flutter install --release
```

#### Android App Bundle (Google Play)

```bash
# Build app bundle
fvm flutter build appbundle --release

# Output: build/app/outputs/bundle/release/app-release.aab
```

#### iOS App (macOS only)

```bash
# Build iOS
fvm flutter build ios --release

# Open Xcode to archive
open ios/Runner.xcworkspace
```

---

## Complete Testing Flow

### 1. Start Backend
```bash
cd backend
go run main.go
```

### 2. Verify Backend
```bash
# In another terminal
curl http://localhost:8080/health
# Should return: ok
```

### 3. Start Frontend
```bash
cd frontend
fvm flutter run
```

### 4. Test User Flow

1. **Register a new user**
   - Open app
   - Click "Register"
   - Enter name, email, password
   - Submit

2. **Create a group**
   - Click "New Group" FAB
   - Enter group name and description
   - Submit

3. **Add members**
   - Open group details
   - Click "Add Member" icon
   - Enter email of existing user
   - Add

4. **Create an expense**
   - In group details, click "Add Expense" FAB
   - Enter title, amount
   - Select members
   - Use "Equal Split" or enter custom amounts
   - Submit

5. **View expense**
   - Tap on expense in group
   - View who paid and who owes

6. **Test logout**
   - Go to Profile
   - Click Logout
   - Verify redirected to login

7. **Test persistence**
   - Close app completely
   - Reopen app
   - Should auto-login

---

## Troubleshooting

### Backend Issues

**Problem:** Database connection failed
```bash
# Check PostgreSQL is running
sudo systemctl status postgresql

# Check credentials in .env
cat .env

# Test database connection
psql -U username -d shared_expenses -h localhost
```

**Problem:** Port 8080 already in use
```bash
# Find process using port
lsof -i :8080

# Kill process
kill -9 <PID>

# Or change port in .env
API_PORT=8081
```

### Frontend Issues

**Problem:** Cannot connect to backend
```bash
# Check backend is running
curl http://localhost:8080/health

# Verify API URL in api_config.dart
# Android emulator: 10.0.2.2:8080
# iOS simulator: localhost:8080
# Physical device: <your-ip>:8080
```

**Problem:** Build errors
```bash
# Clean and rebuild
fvm flutter clean
fvm flutter pub get
fvm flutter run
```

**Problem:** Emulator not starting
```bash
# List emulators
fvm flutter emulators

# Create new emulator in Android Studio
# Tools > Device Manager > Create Device
```

---

## Development Tips

### Hot Reload (Flutter)

While app is running:
- Press `r` - Hot reload (fast, preserves state)
- Press `R` - Hot restart (slower, resets state)
- Press `q` - Quit

### Live Backend Changes

```bash
# Use air for hot reload
go install github.com/cosmtrek/air@latest

# Run with air
air
```

### Debugging

**Flutter:**
```bash
# View logs
fvm flutter logs

# Run in debug mode
fvm flutter run --debug

# Enable verbose logging
fvm flutter run -v
```

**Backend:**
```bash
# Run with verbose logging
go run main.go -v
```

### Code Quality

**Flutter:**
```bash
# Analyze code
fvm flutter analyze

# Format code
fvm flutter format lib/

# Run tests
fvm flutter test
```

**Backend:**
```bash
# Format code
go fmt ./...

# Run tests
go test ./...

# Lint
golangci-lint run
```

---

## Environment-Specific Builds

### Development
- Debug mode
- Local backend (localhost:8080)
- Detailed logging

### Staging
- Release mode
- Staging backend URL
- Reduced logging

### Production
- Release mode
- Production backend URL
- Minimal logging
- Code obfuscation

**Configure in api_config.dart:**
```dart
class ApiConfig {
  static const String baseUrl = String.fromEnvironment(
    'API_URL',
    defaultValue: 'http://10.0.2.2:8080',
  );
}
```

**Build with environment:**
```bash
flutter build apk --release --dart-define=API_URL=https://api.production.com
```

---

## Quick Command Reference

### Backend
```bash
go run main.go                    # Run server
go build                          # Build binary
go test ./...                     # Run tests
go mod tidy                       # Clean dependencies
```

### Frontend
```bash
fvm flutter run                   # Run app
fvm flutter build apk             # Build APK
fvm flutter test                  # Run tests
fvm flutter clean                 # Clean build
fvm flutter pub get               # Get dependencies
fvm flutter doctor                # Check setup
fvm flutter devices               # List devices
```

---

## Next Steps

After successful build and testing:

1. **Customize branding**
   - Update app name in pubspec.yaml
   - Change app icon
   - Update theme colors in main.dart

2. **Add features**
   - Implement expense editing
   - Add settlement calculation
   - Implement data export

3. **Deploy backend**
   - Set up production server
   - Configure HTTPS
   - Set up database backups

4. **Publish app**
   - Android: Google Play Console
   - iOS: App Store Connect

---

## Support

- **Backend docs:** `backend/README.md`
- **Frontend docs:** `frontend/BUILD_INSTRUCTIONS.md`
- **API docs:** `FRONTEND_REQUIREMENTS.md`

---

**Happy Building! 🚀**
