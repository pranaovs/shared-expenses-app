# Flutter Frontend Development Requirements for Shared Expenses Application

## Project Overview
Build a **Flutter mobile application** for a FOSS (GPL v3) expense sharing platform. The application connects to an existing Go backend API and provides an ad-free, privacy-focused alternative to proprietary expense-sharing apps like Splitwise.

---

## Backend API Documentation

### Base URL
- Development: `http://localhost:8080`
- Production: To be configured

### Authentication
All authenticated endpoints require a **JWT Bearer token** in the `Authorization` header:
```
Authorization: Bearer <jwt_token>
```

---

## API Endpoints Reference

### **Authentication Endpoints** (`/auth`)

#### 1. Register New User
- **POST** `/auth/register`
- **Body**:
```json
{
  "name": "John Doe",
  "email": "john@example.com",
  "password": "securepassword"
}
```
- **Response** (200 OK):
```json
{
  "message": "user registered successfully",
  "user_id": "uuid-string"
}
```
- **Errors**: 
  - 400: Validation errors (invalid email, missing fields)
  - 500: User already exists or database error

#### 2. Login
- **POST** `/auth/login`
- **Body**:
```json
{
  "email": "john@example.com",
  "password": "securepassword"
}
```
- **Response** (200 OK):
```json
{
  "message": "login successful",
  "token": "jwt-token-string"
}
```
- **Errors**:
  - 401: Invalid email or password
  - 400: Validation errors

#### 3. Get Current User Profile
- **GET** `/auth/me`
- **Headers**: `Authorization: Bearer <token>`
- **Response** (200 OK):
```json
{
  "user_id": "uuid",
  "name": "John Doe",
  "email": "john@example.com",
  "guest": false,
  "created_at": 1234567890
}
```
- **Errors**:
  - 401: Unauthorized (invalid/missing token)

---

### **User Endpoints** (`/users`)

#### 1. Get User by ID
- **GET** `/users/:id`
- **Headers**: `Authorization: Bearer <token>`
- **Response** (200 OK):
```json
{
  "user_id": "uuid",
  "name": "Jane Smith",
  "email": "jane@example.com",
  "guest": false,
  "created_at": 1234567890
}
```
- **Errors**:
  - 401: Unauthorized
  - 403: Access denied (users not related through any group)
  - 500: Internal server error

#### 2. Search User by Email
- **GET** `/users/search/email/:email`
- **Headers**: `Authorization: Bearer <token>`
- **Response** (200 OK):
```json
{
  "user_id": "uuid",
  "name": "Alice Cooper",
  "email": "alice@example.com",
  "guest": false,
  "created_at": 1234567890
}
```
- **Use Case**: Finding users to add to groups by email
- **Errors**:
  - 401: Unauthorized
  - 400: Invalid email format
  - 500: User not found or database error

---

### **Group Endpoints** (`/groups`)

#### 1. Create Group
- **POST** `/groups/`
- **Headers**: `Authorization: Bearer <token>`
- **Body**:
```json
{
  "name": "Weekend Trip 2024",
  "description": "Expenses for our camping trip"
}
```
- **Response** (200 OK):
```json
{
  "group_id": "uuid"
}
```
- **Note**: Creator is automatically added as the first member and becomes the admin
- **Errors**:
  - 401: Unauthorized
  - 400: Validation errors (missing name)

#### 2. List Groups User is Member Of
- **GET** `/groups/me`
- **Headers**: `Authorization: Bearer <token>`
- **Response** (200 OK):
```json
[
  {
    "group_id": "uuid1",
    "name": "Weekend Trip",
    "description": "Trip expenses",
    "created_by": "uuid",
    "created_at": 1234567890,
    "members": []
  },
  {
    "group_id": "uuid2",
    "name": "Roommates",
    "description": "Apartment expenses",
    "created_by": "uuid",
    "created_at": 1234567891,
    "members": []
  }
]
```
- **Note**: Includes groups created by user and groups user was added to
- **Errors**:
  - 401: Unauthorized

#### 3. List Groups User is Admin Of
- **GET** `/groups/admin`
- **Headers**: `Authorization: Bearer <token>`
- **Response** (200 OK):
```json
[
  {
    "group_id": "uuid1",
    "name": "Weekend Trip",
    "description": "Trip expenses",
    "created_by": "uuid",
    "created_at": 1234567890,
    "members": []
  }
]
```
- **Use Case**: Show groups where user has admin privileges (can add/remove members)
- **Errors**:
  - 401: Unauthorized

#### 4. Get Group Details
- **GET** `/groups/:id`
- **Headers**: `Authorization: Bearer <token>`
- **Response** (200 OK):
```json
{
  "group_id": "uuid",
  "name": "Weekend Trip",
  "description": "Camping trip expenses",
  "created_by": "creator-uuid",
  "created_at": 1234567890,
  "members": [
    {
      "user_id": "uuid1",
      "name": "John Doe",
      "email": "john@example.com",
      "guest": false,
      "joined_at": 1234567890
    },
    {
      "user_id": "uuid2",
      "name": "Jane Smith",
      "email": "jane@example.com",
      "guest": false,
      "joined_at": 1234567891
    }
  ]
}
```
- **Errors**:
  - 401: Unauthorized
  - 403: Access denied (not a member)
  - 404: Group not found

#### 5. Add Members to Group
- **POST** `/groups/:id/members`
- **Headers**: `Authorization: Bearer <token>`
- **Body**:
```json
{
  "user_ids": ["uuid1", "uuid2", "uuid3"]
}
```
- **Response** (200 OK):
```json
{
  "message": "members added successfully",
  "added_members": ["uuid1", "uuid2"]
}
```
- **Note**: Only group admin (creator) can add members. Invalid user IDs are silently skipped.
- **Errors**:
  - 401: Unauthorized
  - 403: Only group admin can add members
  - 404: Group not found
  - 400: No valid user IDs

#### 6. Remove Members from Group
- **DELETE** `/groups/:id/members`
- **Headers**: `Authorization: Bearer <token>`
- **Body**:
```json
{
  "user_ids": ["uuid1", "uuid2"]
}
```
- **Response** (200 OK):
```json
{
  "message": "members removed",
  "removed_members": ["uuid1", "uuid2"]
}
```
- **Note**: Only group admin can remove members. Cannot remove the admin themselves.
- **Errors**:
  - 401: Unauthorized
  - 403: Only group admin can remove members
  - 404: Group not found
  - 400: Cannot remove group admin

---

### **Expense Endpoints** (`/expenses`)

#### 1. Create Expense
- **POST** `/expenses/`
- **Headers**: `Authorization: Bearer <token>`
- **Body**:
```json
{
  "group_id": "group-uuid",
  "title": "Groceries",
  "description": "Weekly grocery shopping",
  "amount": 120.50,
  "is_incomplete_amount": false,
  "is_incomplete_split": false,
  "latitude": 37.7749,
  "longitude": -122.4194,
  "splits": [
    {
      "user_id": "user1-uuid",
      "amount": 120.50,
      "is_paid": true
    },
    {
      "user_id": "user2-uuid",
      "amount": 60.25,
      "is_paid": false
    },
    {
      "user_id": "user3-uuid",
      "amount": 60.25,
      "is_paid": false
    }
  ]
}
```
- **Response** (200 OK):
```json
{
  "expense_id": "expense-uuid"
}
```

**Splits Explanation**:
- `is_paid: true` = User **paid** this amount (lender/contributor)
- `is_paid: false` = User **owes** this amount (borrower)
- The sum of all `is_paid: true` amounts should equal the total expense `amount`
- The sum of all `is_paid: false` amounts should equal the total expense `amount`
- Same user can appear multiple times with different `is_paid` values

**Incomplete Flags**:
- `is_incomplete_amount: true` = Amount validation is skipped (for partial entries)
- `is_incomplete_split: true` = Split validation is skipped (for incomplete data)

**Validation**:
- User must be a member of the group
- All users in splits must be members of the group
- If incomplete flags are false, paid/owed totals must match expense amount (within tolerance of 0.01)

**Errors**:
  - 401: Unauthorized
  - 403: User not a member of group
  - 400: Validation errors (no splits, split user not in group, amount mismatch)

#### 2. Get Expense by ID
- **GET** `/expenses/:id`
- **Headers**: `Authorization: Bearer <token>`
- **Response** (200 OK):
```json
{
  "expense_id": "uuid",
  "group_id": "group-uuid",
  "added_by": "user-uuid",
  "title": "Groceries",
  "description": "Weekly shopping",
  "created_at": 1234567890,
  "amount": 120.50,
  "is_incomplete_amount": false,
  "is_incomplete_split": false,
  "latitude": 37.7749,
  "longitude": -122.4194,
  "splits": [
    {
      "user_id": "user1-uuid",
      "amount": 120.50,
      "is_paid": true
    },
    {
      "user_id": "user2-uuid",
      "amount": 60.25,
      "is_paid": false
    }
  ]
}
```
- **Errors**:
  - 401: Unauthorized
  - 403: Access denied (not a group member)
  - 404: Expense not found

#### 3. Update Expense
- **PUT** `/expenses/:id`
- **Headers**: `Authorization: Bearer <token>`
- **Body**: Same structure as create expense (without `expense_id`)
```json
{
  "group_id": "group-uuid",
  "title": "Updated title",
  "description": "Updated description",
  "amount": 150.00,
  "is_incomplete_amount": false,
  "is_incomplete_split": false,
  "latitude": 37.7749,
  "longitude": -122.4194,
  "splits": [...]
}
```
- **Response** (200 OK):
```json
{
  "message": "expense updated"
}
```
- **Authorization**: Only the user who added the expense OR the group admin can update
- **Note**: All splits are replaced (not merged)
- **Errors**:
  - 401: Unauthorized
  - 403: Not authorized (not expense adder or group admin)
  - 404: Expense not found
  - 400: Validation errors

#### 4. Delete Expense
- **DELETE** `/expenses/:id`
- **Headers**: `Authorization: Bearer <token>`
- **Response** (200 OK):
```json
{
  "message": "expense deleted"
}
```
- **Authorization**: Only the user who added the expense OR the group admin can delete
- **Errors**:
  - 401: Unauthorized
  - 403: Not authorized
  - 404: Expense not found

---

## Database Schema Overview

### Users Table
- `user_id` (UUID, Primary Key)
- `user_name` (TEXT)
- `email` (TEXT, Unique)
- `password_hash` (TEXT, nullable)
- `is_guest` (BOOLEAN, default false)
- `created_at` (TIMESTAMPTZ)

### Groups Table
- `group_id` (UUID, Primary Key)
- `group_name` (TEXT)
- `description` (TEXT, nullable)
- `created_by` (UUID, Foreign Key → users)
- `created_at` (TIMESTAMPTZ)

### Group Members Table
- `user_id` (UUID, Foreign Key → users)
- `group_id` (UUID, Foreign Key → groups)
- `joined_at` (TIMESTAMPTZ)
- Primary Key: (user_id, group_id)

### Expenses Table
- `expense_id` (UUID, Primary Key)
- `group_id` (UUID, Foreign Key → groups)
- `added_by` (UUID, Foreign Key → users)
- `title` (TEXT)
- `description` (TEXT, nullable)
- `created_at` (TIMESTAMPTZ)
- `amount` (DOUBLE PRECISION)
- `is_incomplete_amount` (BOOLEAN)
- `is_incomplete_split` (BOOLEAN)
- `latitude` (DOUBLE PRECISION, nullable)
- `longitude` (DOUBLE PRECISION, nullable)

### Expense Splits Table
- `expense_id` (UUID, Foreign Key → expenses)
- `user_id` (UUID, Foreign Key → users)
- `amount` (DOUBLE PRECISION)
- `is_paid` (BOOLEAN) - true = lender, false = borrower
- Primary Key: (expense_id, user_id, is_paid)

---

## Flutter Frontend Requirements

### 1. **Technology Stack**
- **Framework**: Flutter (latest stable version)
- **Language**: Dart
- **State Management**: Provider, Riverpod, or Bloc (choose based on expertise)
- **HTTP Client**: `http` or `dio` package
- **Local Storage**: `shared_preferences` for token storage, `hive` or `sqflite` for offline capability
- **Navigation**: Named routes with `go_router` or similar

### 2. **Core Features to Implement**

#### **Authentication Flow**
1. **Splash Screen**
   - Check if user has valid JWT token
   - Auto-login if token exists and is valid
   - Navigate to login/register if not authenticated

2. **Login Screen**
   - Email and password input fields
   - Form validation (email format, required fields)
   - Login button → POST `/auth/login`
   - Store JWT token securely in local storage
   - Navigate to home screen on success
   - Show error messages for failed login
   - "Register" link to navigate to registration

3. **Registration Screen**
   - Name, email, password input fields
   - Password confirmation field
   - Form validation
   - Register button → POST `/auth/register`
   - Auto-login after successful registration
   - Show success/error messages
   - "Login" link to return to login screen

4. **Profile Screen**
   - Display user information from GET `/auth/me`
   - Show name, email, account creation date
   - Logout button (clear token and return to login)
   - Optional: Edit profile (future feature)

#### **Group Management**

1. **Groups List Screen (Home)**
   - Fetch and display groups from GET `/groups/me`
   - Show group name, description, member count
   - Pull-to-refresh functionality
   - Floating action button to create new group
   - Tap group to navigate to group details
   - Bottom navigation: Groups, Profile tabs
   - Optional: Show groups user is admin of separately (GET `/groups/admin`)

2. **Create Group Screen**
   - Group name input (required)
   - Description input (optional)
   - Create button → POST `/groups/`
   - Navigate back to groups list on success
   - Show success/error feedback

3. **Group Details Screen**
   - Fetch group details from GET `/groups/:id`
   - Display group name, description, creation date
   - **Members Section**:
     - List all members with name, email
     - Show member join date
     - If user is admin: Show add/remove member buttons
   - **Add Members Dialog** (admin only):
     - Search user by email (GET `/users/search/email/:email`)
     - Add found user to group (POST `/groups/:id/members`)
     - Show feedback messages
   - **Remove Member** (admin only):
     - Long-press or swipe member to remove
     - Confirm deletion dialog
     - Send DELETE `/groups/:id/members`
   - **Expenses List**:
     - Show recent expenses for this group (fetch from backend or cache)
     - Tap expense to view details
   - **Floating Action Button**: Add new expense

#### **Expense Management**

1. **Create Expense Screen**
   - **Basic Information**:
     - Title (required)
     - Description (optional)
     - Amount (required, numeric input)
     - Date picker (default: today)
   - **Location** (optional):
     - Use GPS to capture latitude/longitude
     - Show map preview if available
   - **Incomplete Flags**:
     - Checkbox: "Amount is incomplete"
     - Checkbox: "Split is incomplete"
   - **Split Configuration**:
     - Show list of group members
     - For each member, allow setting:
       - Paid amount (is_paid: true)
       - Owed amount (is_paid: false)
     - **Quick Actions**:
       - "Equal split" button: Divide amount equally among selected members
       - "Split by percentage" option
       - Auto-calculate remaining amount
     - Validate: Sum of paid amounts = expense amount
     - Validate: Sum of owed amounts = expense amount
   - **Submit Button**: POST `/expenses/`
   - Show success/error feedback
   - Navigate back to group details on success

2. **Expense Details Screen**
   - Fetch expense from GET `/expenses/:id`
   - Display:
     - Title, description, amount
     - Creation date and time
     - Added by (user name)
     - Location on map if available
   - **Splits Breakdown**:
     - Section: "Who Paid"
       - List users with `is_paid: true` and amounts
     - Section: "Who Owes"
       - List users with `is_paid: false` and amounts
   - **Edit Button** (if user is expense creator or group admin):
     - Navigate to edit expense screen
   - **Delete Button** (if user is expense creator or group admin):
     - Show confirmation dialog
     - DELETE `/expenses/:id`
     - Navigate back on success

3. **Edit Expense Screen**
   - Same UI as create expense screen
   - Pre-populate fields with existing expense data
   - PUT `/expenses/:id` on save

#### **Expense Calculation & Settlement**

**Note**: Settlement calculation is not yet implemented in the backend. The frontend should be designed to support future implementation.

1. **Settlement View** (Future Feature)
   - Calculate who owes whom based on all expenses
   - Display simplified debt graph
   - Show minimum transactions needed to settle all debts
   - Mark expenses as "settled"

For now, implement a basic calculation locally:
- For a group, fetch all expenses
- Calculate net balance for each user (total paid - total owed)
- Display who owes whom

### 3. **UI/UX Design Guidelines**

#### **Design Principles**
- **Clean & Minimal**: No ads, no clutter
- **Material Design 3**: Use Flutter's Material 3 theming
- **Color Scheme**: 
  - Primary: Green/Blue tones (trust, money)
  - Accent: Orange/Yellow (warmth, friendliness)
  - Background: Light mode default, support dark mode
- **Accessibility**: 
  - High contrast ratios
  - Font scaling support
  - Screen reader compatibility

#### **Key UI Components**

1. **Group Card**
   - Card with shadow elevation
   - Group name (headline text)
   - Description (body text, 2 lines max with ellipsis)
   - Member count badge
   - Tap to open group details

2. **Expense List Item**
   - Title and amount (bold)
   - Date and added by (caption)
   - Category icon (if implemented)
   - Swipe actions: Edit, Delete (if authorized)

3. **Member Chip**
   - User avatar (initials if no image)
   - Name
   - Badge for admin/creator

4. **Split Input Widget**
   - User name
   - Two input fields side-by-side:
     - "Paid" amount (green border)
     - "Owes" amount (orange border)
   - Checkbox to include/exclude user from split

5. **Empty States**
   - Friendly illustrations for empty groups, no expenses
   - Call-to-action buttons

6. **Error Handling**
   - Toast messages for network errors
   - Retry buttons
   - Offline mode indicators

### 4. **State Management Architecture**

#### **Recommended Pattern**: Provider or Riverpod

**State Structure**:
```dart
// Authentication State
class AuthState {
  String? token;
  User? currentUser;
  bool isAuthenticated;
  bool isLoading;
}

// Groups State
class GroupsState {
  List<Group> userGroups;
  List<Group> adminGroups;
  Group? selectedGroup;
  bool isLoading;
  String? error;
}

// Expenses State
class ExpensesState {
  Map<String, List<Expense>> expensesByGroup;
  Expense? selectedExpense;
  bool isLoading;
  String? error;
}
```

**Services**:
```dart
// API Service
class ApiService {
  Future<String> login(String email, String password);
  Future<String> register(String name, String email, String password);
  Future<User> getCurrentUser();
  Future<List<Group>> getMyGroups();
  Future<Group> getGroup(String groupId);
  Future<String> createGroup(String name, String description);
  Future<void> addMembers(String groupId, List<String> userIds);
  Future<void> removeMembers(String groupId, List<String> userIds);
  Future<User> searchUserByEmail(String email);
  Future<String> createExpense(ExpenseRequest expense);
  Future<Expense> getExpense(String expenseId);
  Future<void> updateExpense(String expenseId, ExpenseRequest expense);
  Future<void> deleteExpense(String expenseId);
}

// Storage Service
class StorageService {
  Future<void> saveToken(String token);
  Future<String?> getToken();
  Future<void> deleteToken();
}
```

### 5. **Data Models**

```dart
class User {
  final String userId;
  final String name;
  final String email;
  final bool guest;
  final int createdAt;
}

class Group {
  final String groupId;
  final String name;
  final String? description;
  final String createdBy;
  final int createdAt;
  final List<GroupUser> members;
}

class GroupUser {
  final String userId;
  final String name;
  final String email;
  final bool guest;
  final int joinedAt;
}

class Expense {
  final String expenseId;
  final String groupId;
  final String addedBy;
  final String title;
  final String? description;
  final int createdAt;
  final double amount;
  final bool isIncompleteAmount;
  final bool isIncompleteSplit;
  final double? latitude;
  final double? longitude;
  final List<ExpenseSplit> splits;
}

class ExpenseSplit {
  final String userId;
  final double amount;
  final bool isPaid; // true = paid, false = owes
}
```

### 6. **Error Handling**

**HTTP Error Codes**:
- **401 Unauthorized**: Clear token, redirect to login
- **403 Forbidden**: Show "Access Denied" message
- **404 Not Found**: Show "Resource not found"
- **400 Bad Request**: Show validation error messages
- **500 Server Error**: Show "Server error, try again later"

**Network Errors**:
- Catch `SocketException` for no internet
- Show offline banner
- Queue actions for retry when online (optional)

### 7. **Testing Requirements**

1. **Unit Tests**:
   - Test API service methods
   - Test data model serialization/deserialization
   - Test state management logic

2. **Widget Tests**:
   - Test authentication screens
   - Test group list and details
   - Test expense creation form validation

3. **Integration Tests**:
   - Test complete user flows (register → create group → add expense)

### 8. **Security Considerations**

1. **Token Storage**: Use `flutter_secure_storage` for JWT tokens (encrypted storage)
2. **HTTPS Only**: Ensure all API calls use HTTPS in production
3. **Input Validation**: Validate all user inputs before sending to API
4. **Password Handling**: Never log or store passwords in plain text
5. **Logout**: Clear all local data and tokens on logout

### 9. **Offline Support (Optional Future Enhancement)**

1. Use `sqflite` to cache:
   - User groups
   - Recent expenses
   - User profiles
2. Sync data when online
3. Show "offline" indicator
4. Queue create/update/delete operations for later sync

### 10. **Performance Optimization**

1. **Pagination**: If expense lists grow large, implement pagination
2. **Image Optimization**: Compress expense images before upload (future feature)
3. **Lazy Loading**: Load group details only when needed
4. **Caching**: Cache API responses appropriately
5. **Debouncing**: Debounce search inputs (email search)

---

## Screen Flow Diagram

```
Splash Screen
    ↓
    ├→ (No Token) → Login Screen ⇄ Register Screen
    ↓                                     ↓
    └→ (Has Token) → Groups List Screen ←┘
                            ↓
                            ├→ Create Group Screen
                            ↓
                            └→ Group Details Screen
                                    ↓
                                    ├→ Add Member Dialog
                                    ├→ Remove Member (admin only)
                                    ↓
                                    └→ Create Expense Screen
                                            ↓
                                            └→ Expense Details Screen
                                                    ↓
                                                    ├→ Edit Expense Screen
                                                    └→ Delete Expense
```

---

## Implementation Priority

### **Phase 1: Core Authentication & Navigation** (Week 1)
- [ ] Setup Flutter project
- [ ] Implement API service layer
- [ ] Create data models
- [ ] Build authentication screens (login, register)
- [ ] Implement token storage
- [ ] Setup navigation

### **Phase 2: Group Management** (Week 2)
- [ ] Groups list screen
- [ ] Create group screen
- [ ] Group details screen
- [ ] Add/remove members (admin only)
- [ ] User search by email

### **Phase 3: Expense Management** (Week 3)
- [ ] Create expense screen with split logic
- [ ] Expense details screen
- [ ] Edit expense
- [ ] Delete expense
- [ ] Validate split amounts

### **Phase 4: Polish & Testing** (Week 4)
- [ ] Error handling
- [ ] Loading states
- [ ] Empty states
- [ ] Pull-to-refresh
- [ ] Unit and widget tests
- [ ] Dark mode support
- [ ] Accessibility improvements

### **Phase 5: Future Enhancements**
- [ ] Settlement calculation
- [ ] Expense filtering
- [ ] Expense categories
- [ ] Export data (CSV, JSON, XML)
- [ ] Offline support
- [ ] Expense images
- [ ] Notifications
- [ ] Guest user support

---

## Additional Notes

### **Missing Backend Features** (Not Yet Implemented)
1. **List expenses for a group**: You'll need to request backend team to add `GET /groups/:id/expenses`
2. **Expense filtering**: No backend endpoints for filtering by date, amount, user
3. **Settlement calculation**: No backend logic for calculating net balances
4. **Guest users**: Backend has `is_guest` field but no guest-specific endpoints
5. **Expense categories**: Not in current schema
6. **Expense images**: No image upload endpoints
7. **Expense edit history**: Not tracked in database

### **Workarounds for Frontend**
- **Listing group expenses**: Fetch all groups' expenses and filter locally
- **Settlement calculation**: Implement client-side calculation using fetched expenses
- **Categories**: Add as local-only tags until backend supports them

### **Backend API Gaps to Address**
Ask backend team to implement:
1. `GET /groups/:id/expenses` - List all expenses for a group
2. `GET /expenses?group_id=X&start_date=Y&end_date=Z` - Filter expenses
3. `GET /groups/:id/balance` - Get net balance per user
4. `GET /groups/:id/settlements` - Get optimal settlement transactions

---

## Environment Configuration

### **Development**
```dart
const String API_BASE_URL = 'http://10.0.2.2:8080'; // Android emulator
// const String API_BASE_URL = 'http://localhost:8080'; // iOS simulator
```

### **Production**
```dart
const String API_BASE_URL = 'https://api.yourapp.com';
```

Use `flutter_dotenv` to manage environment variables.

---

## Dependencies Suggestion

```yaml
dependencies:
  flutter:
    sdk: flutter
  
  # State Management
  provider: ^6.1.1
  # OR
  flutter_riverpod: ^2.4.9
  
  # HTTP Client
  dio: ^5.4.0
  
  # Local Storage
  shared_preferences: ^2.2.2
  flutter_secure_storage: ^9.0.0
  
  # Navigation
  go_router: ^13.0.0
  
  # UI Components
  flutter_svg: ^2.0.9
  cached_network_image: ^3.3.1
  
  # Date/Time
  intl: ^0.19.0
  
  # Location (for GPS)
  geolocator: ^11.0.0
  
  # Optional: Maps
  flutter_map: ^6.1.0
  
dev_dependencies:
  flutter_test:
    sdk: flutter
  mockito: ^5.4.4
  build_runner: ^2.4.8
```

---

## Success Criteria

Your Flutter frontend should:
1. ✅ Allow users to register and login
2. ✅ Display all groups the user is a member of
3. ✅ Allow creating new groups
4. ✅ Show group details with all members
5. ✅ Allow admins to add/remove members
6. ✅ Allow creating expenses with uneven splits
7. ✅ Show expense details with split breakdown
8. ✅ Allow authorized users to edit/delete expenses
9. ✅ Handle all API errors gracefully
10. ✅ Persist authentication across app restarts
11. ✅ Support both light and dark themes
12. ✅ Follow Material Design 3 guidelines
13. ✅ Be responsive on various screen sizes

---

## Questions to Clarify with Backend Team

1. Is there a health check endpoint to verify API availability? (Yes: GET `/health`)
2. What's the JWT token expiration time? Should we implement refresh tokens?
3. How do we handle concurrent expense edits?
4. Is there a rate limit on API requests?
5. Will you add `GET /groups/:id/expenses` endpoint?
6. Can we get expense edit history in the future?
7. Should we support expense attachments (images)?

---

## License & Contribution

This project is licensed under **GPL v3**. All contributions must maintain:
- No advertisements
- No tracking
- Open source
- User privacy first

---

**Good luck with the frontend development!** 🚀

For any questions, refer to the backend source code at `/home/personal/repos/shared-expenses-app/backend/`.
