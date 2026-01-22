// Package app provides implementation for app
//
// File: dependency.go
// Description: implementation for app
// Author: Bayarsaikhan Otgonbayar, CTO
// Company: Gerege Core Team
// Created: 2025-02-20
// Last Updated: 2025-02-20
/*
Package app нь application-ийн dependency injection (DI) container-ийг тодорхойлно.

Энэ package нь Clean Architecture-ийн "Composition Root" үүрэг гүйцэтгэнэ:
  - Бүх repository-уудыг үүсгэнэ
  - Бүх service-уудыг үүсгэж, шаардлагатай dependency-уудыг inject хийнэ
  - Handler-уудад хэрэгтэй бүх зүйлийг Dependencies struct-ээр дамжуулна

Dependency Graph:

	┌─────────────────────────────────────────────────────┐
	│                   Dependencies                       │
	├─────────────────────────────────────────────────────┤
	│  ┌─────────┐  ┌─────────┐  ┌─────────┐  ┌────────┐  │
	│  │   DB    │  │  Config │  │  Logger │  │ SSO    │  │
	│  └────┬────┘  └────┬────┘  └────┬────┘  └────────┘  │
	│       │            │            │                    │
	│       ▼            ▼            ▼                    │
	│  ┌─────────────────────────────────────────────┐    │
	│  │              RepoContainer                   │    │
	│  │  (User, Role, Organization, Module, etc.)   │    │
	│  └────────────────────┬────────────────────────┘    │
	│                       │                              │
	│                       ▼                              │
	│  ┌─────────────────────────────────────────────┐    │
	│  │             ServiceContainer                 │    │
	│  │  (UserService, RoleService, etc.)           │    │
	│  └─────────────────────────────────────────────┘    │
	└─────────────────────────────────────────────────────┘

Ашиглалт:

	deps := app.NewDependencies(db, cfg, logger, authCache)
	router.MapV1(app, deps)
*/
package app

import (
	"time"

	"git.gerege.mn/backend-packages/config"     // Application configuration
	"git.gerege.mn/backend-packages/sso-client" // SSO client
	"templatev25/internal/auth"                 // Permission cache
	localconfig "templatev25/internal/config"   // Local auth config
	"templatev25/internal/repository"           // Data access layer
	"templatev25/internal/service"              // Business logic layer

	"github.com/redis/go-redis/v9" // Redis client
	"go.uber.org/zap"              // Structured logging
	"gorm.io/gorm"                 // ORM
)

// ============================================================
// DEPENDENCIES STRUCT
// ============================================================

// Dependencies нь application-ийн бүх dependency-уудыг агуулна.
// Энэ struct-ийг router, handler-уудад дамжуулна.
//
// Fields:
//   - DB: GORM database connection (PostgreSQL)
//   - Log: Zap structured logger
//   - Cfg: Application configuration
//   - AuthCache: Session cache (LRU)
//   - SSO: SSO HTTP client (authentication)
//   - Repo: Repository container (data access)
//   - Service: Service container (business logic)
type Dependencies struct {
	// DB нь GORM database connection.
	// Бүх repository-ууд энэ connection-ийг ашиглана.
	DB *gorm.DB

	// Log нь structured logger.
	// Бүх log бичилтэд ашиглана (info, error, debug).
	Log *zap.Logger

	// Cfg нь application configuration.
	// Environment variables, .env файлаас уншсан тохиргоо.
	Cfg *config.Config

	// AuthCache нь session cache.
	// SSO-оос ирсэн session-уудыг LRU cache-д хадгална.
	// Дахин SSO руу request илгээхгүйгээр session validate хийнэ.
	AuthCache *ssoclient.Cache

	// SSO нь SSO HTTP client.
	// OAuth2 flow, session validation зэрэгт ашиглана.
	SSO *ssoclient.SSOClient

	// PermCache нь permission cache.
	// Permission шалгахад ашиглана.
	// auth.RequirePermission middleware-д дамжуулна.
	PermCache *auth.PermissionCache

	// Repo нь бүх repository-уудыг агуулна.
	// Database CRUD operations.
	Repo *RepoContainer

	// Service нь бүх service-уудыг агуулна.
	// Business logic, validation, external API calls.
	Service *ServiceContainer
}

// ============================================================
// REPOSITORY CONTAINER
// ============================================================

// RepoContainer нь бүх repository instance-уудыг агуулна.
// Repository нь database-тай шууд харьцах давхарга (Data Access Layer).
//
// Naming convention:
//   - Interface: repository.UserRepository
//   - Implementation: repository.NewUserRepository(db)
//
// Repository responsibilities:
//   - CRUD operations (Create, Read, Update, Delete)
//   - Pagination, filtering, sorting
//   - Transaction management
//   - SQL query optimization
type RepoContainer struct {
	// ============================================================
	// USER & AUTH REPOSITORIES
	// ============================================================

	// User нь хэрэглэгчийн CRUD operations.
	// Table: users
	User repository.UserRepository

	// UserRole нь хэрэглэгч-эрхийн холбоос.
	// Table: user_roles (many-to-many)
	UserRole repository.UserRoleRepository

	// Auth нь local authentication CRUD operations.
	// Tables: user_credentials, user_mfa_totp, sessions, login_history, etc.
	Auth repository.AuthRepository

	// Registration нь registration-related CRUD operations.
	// Tables: email_verification_tokens, password_reset_tokens
	Registration repository.RegistrationRepository

	// ============================================================
	// SYSTEM & MODULE REPOSITORIES
	// ============================================================

	// System нь системийн CRUD operations.
	// Table: systems (app groups)
	System repository.SystemRepository

	// Module нь модулийн CRUD operations.
	// Table: modules (menu items)
	Module repository.ModuleRepository


	// Menu нь цэсний CRUD operations.
	// Table: menus
	Menu repository.MenuRepository

	// ============================================================
	// PERMISSION & ROLE REPOSITORIES
	// ============================================================

	// Permission нь зөвшөөрлийн CRUD operations.
	// Table: permissions
	Permission repository.PermissionRepository

	// Action нь action-ийн CRUD operations.
	// Table: actions
	Action repository.ActionRepository

	// Role нь эрхийн CRUD operations.
	// Table: roles
	Role repository.RoleRepository

	// ============================================================
	// ORGANIZATION REPOSITORIES
	// ============================================================

	// Organization нь байгууллагын CRUD operations.
	// Table: organizations
	Organization repository.OrganizationRepository

	// OrganizationType нь байгууллагын төрлийн CRUD operations.
	// Table: organization_types
	OrganizationType repository.OrganizationTypeRepository

	// OrgUser нь байгууллага-хэрэглэгчийн холбоос.
	// Table: organization_users (many-to-many)
	OrgUser repository.OrgUserRepository

	// ============================================================
	// TERMINAL & PLATFORM REPOSITORIES
	// ============================================================

	// Terminal нь терминалын CRUD operations.
	// Table: terminals
	Terminal repository.TerminalRepository

	// AppServiceIcon нь app service icon-ийн CRUD operations.
	// Table: app_service_icons
	AppServiceIcon repository.AppServiceIconRepository

	// AppServiceIconGroup нь app service icon group-ийн CRUD operations.
	// Table: app_service_icon_groups
	AppServiceIconGroup repository.AppServiceIconGroupRepository


	// ============================================================
	// CONTENT REPOSITORIES
	// ============================================================

	// PublicFile нь нийтийн файлын CRUD operations.
	// Table: public_files
	PublicFile repository.PublicFileRepository

	// Notification нь мэдэгдлийн CRUD operations.
	// Table: notifications
	Notification repository.NotificationRepository

	// News нь мэдээний CRUD operations.
	// Table: news
	News repository.NewsRepository

	// ChatItem нь chat item-ийн CRUD operations.
	// Table: chat_items
	ChatItem repository.ChatItemRepository

	// APILog нь API log-ийн CRUD operations.
	// Table: logs
	APILog repository.APILogRepository
}

// ============================================================
// SERVICE CONTAINER
// ============================================================

// ServiceContainer нь бүх service instance-уудыг агуулна.
// Service нь business logic давхарга (Business Layer).
//
// Service responsibilities:
//   - Business rules, validation
//   - Data transformation (DTO ↔ Domain)
//   - External API calls
//   - Transaction orchestration
//   - Error handling, logging
type ServiceContainer struct {
	// ============================================================
	// USER & AUTH SERVICES
	// ============================================================

	// User нь хэрэглэгчийн business logic.
	// - Profile management
	// - SSO integration
	// - Organization membership
	User *service.UserService

	// UserRole нь хэрэглэгч-эрхийн business logic.
	// - Role assignment
	// - Permission checking
	UserRole service.UserRoleService

	// Auth нь local authentication service.
	// - Login, MFA, password management
	// - Session management
	Auth *service.AuthService

	// SessionStore нь Redis session storage.
	// - Session CRUD
	// - MFA token storage
	SessionStore service.SessionStore

	// Registration нь user registration service.
	// - User signup
	// - Email verification
	// - Password reset
	Registration *service.RegistrationService

	// ============================================================
	// SYSTEM & MODULE SERVICES
	// ============================================================

	// System нь системийн business logic.
	// - System CRUD
	// - System-module relations
	System service.SystemService

	// Module нь модулийн business logic.
	// - Menu management
	// - Access control
	Module service.ModuleService


	// Menu нь цэсний business logic.
	// - Menu CRUD
	// - Hierarchical menu structure
	Menu service.MenuService

	// ============================================================
	// PERMISSION & ROLE SERVICES
	// ============================================================

	// Permission нь зөвшөөрлийн business logic.
	// - Permission CRUD
	// - Permission checking
	Permission *service.PermissionService

	// Action нь action-ийн business logic.
	// - Action CRUD (Permission-тэй ижил логик)
	Action *service.ActionService

	// Role нь эрхийн business logic.
	// - Role CRUD
	// - Role-permission assignment
	Role *service.RoleService

	// ============================================================
	// ORGANIZATION SERVICES
	// ============================================================

	// Organization нь байгууллагын business logic.
	// - Organization CRUD
	// - Core system integration
	Organization *service.OrganizationService

	// OrganizationType нь байгууллагын төрлийн business logic.
	OrganizationType *service.OrganizationTypeService

	// OrgUser нь байгууллага-хэрэглэгчийн business logic.
	// - User membership management
	// - Organization switching
	OrgUser *service.OrgUserService

	// ============================================================
	// TERMINAL & PLATFORM SERVICES
	// ============================================================

	// Terminal нь терминалын business logic.
	Terminal *service.TerminalService

	// AppServiceIcon нь app service icon-ийн business logic.
	AppServiceIcon *service.AppServiceIconService

	// AppServiceGroup нь app service icon group-ийн business logic.
	AppServiceGroup *service.AppServiceIconGroup


	// ============================================================
	// CONTENT SERVICES
	// ============================================================

	// PublicFile нь нийтийн файлын business logic.
	// - File upload/download
	// - Access control
	PublicFile *service.PublicFileService

	// Notification нь мэдэгдлийн business logic.
	// - Send notifications
	// - Mark as read
	Notification *service.NotificationService

	// News нь мэдээний business logic.
	News *service.NewsService

	// ChatItem нь chat item-ийн business logic.
	ChatItem *service.ChatItemService

	// APILog нь API log-ийн business logic.
	// - API log listing with pagination
	APILog service.APILogService

	// ============================================================
	// EXTERNAL INTEGRATION SERVICES
	// ============================================================

	// Verify нь баталгаажуулалтын service.
	// - ХУР (XYP) integration
	// - Passport verification
	Verify *service.VerifyService

	// Meet нь видео хурлын service.
	// - Video conference room management
	Meet *service.MeetService

	// Tpay нь терминал төлбөрийн service.
	// - Payment processing
	// - Card management
	Tpay *service.TpayService
}

// ============================================================
// CONSTRUCTOR FUNCTION
// ============================================================

// NewDependencies нь бүх dependency-уудыг үүсгэж, холбоно.
// Энэ функц нь Composition Root болж ажиллана.
//
// Parameters:
//   - db: GORM database connection
//   - cfg: Application configuration
//   - log: Zap structured logger
//   - authCache: Session cache instance
//
// Returns:
//   - *Dependencies: Бүх dependency-уудыг агуулсан struct
//
// Dependency creation order:
//  1. Repositories (database layer)
//  2. Services (business layer, repositories-ээс хамаарна)
//  3. SSO client (auth layer)
//  4. Final Dependencies struct
func NewDependencies(db *gorm.DB, cfg *config.Config, log *zap.Logger, authCache *ssoclient.Cache) *Dependencies {

	// ============================================================
	// STEP 1: Create all repositories
	// ============================================================
	// Repository-ууд нь database connection-оос хамаарна.
	// Зарим repository-ууд config-оос нэмэлт тохиргоо авна.
	repo := &RepoContainer{
		// User & Auth
		User:         repository.NewUserRepository(db),
		UserRole:     repository.NewUserRoleRepository(db),
		Auth:         repository.NewAuthRepository(db),
		Registration: repository.NewRegistrationRepository(db),

		// System & Module
		System: repository.NewSystemRepository(db),
		Module: repository.NewModuleRepository(db, cfg), // config: table prefix
		Menu:   repository.NewMenuRepository(db, cfg),   // config: schema name

		// Permission & Role
		Permission: repository.NewPermissionRepository(db),
		Action:     repository.NewActionRepository(db),
		Role:       repository.NewRoleRepository(db),

		// Organization
		Organization:     repository.NewOrganizationRepository(db),
		OrganizationType: repository.NewOrganizationTypeRepository(db),
		OrgUser:          repository.NewOrgUserRepository(db, cfg), // config: external URLs

		// Terminal & Platform
		Terminal:            repository.NewTerminalRepository(db),
		AppServiceIcon:      repository.NewAppServiceIconRepository(db),
		AppServiceIconGroup: repository.NewAppServiceIconGroupRepository(db),

		// Content
		PublicFile:   repository.NewPublicFileRepository(db),
		Notification: repository.NewNotificationRepository(db),
		News:         repository.NewNewsRepository(db),
		ChatItem:     repository.NewChatItemRepository(db),

		// Logging
		APILog: repository.NewAPILogRepository(db),
	}

	// ============================================================
	// STEP 2: Create all services
	// ============================================================
	// Service-ууд нь repository-уудаас хамаарна.
	// Зарим service-ууд config, logger, бусад repository-уудыг авна.
	
	// Permission service эхлээд үүсгэх (Action service-д хэрэгтэй)
	permissionSvc := service.NewPermissionService(repo.Permission, log)
	
	svc := &ServiceContainer{
		// User & Auth
		User:     service.NewUserService(repo.User, cfg, log), // External API calls
		UserRole: service.NewUserRoleService(repo.UserRole),

		// System & Module
		System: service.NewSystemService(repo.System, log),
		Module: service.NewModuleService(repo.Module),
		Menu:   service.NewMenuService(repo.Menu),

		// Permission & Role
		Permission: permissionSvc,
		Action:     service.NewActionService(repo.Action, log),
		Role:       service.NewRoleService(repo.Role, log),

		// Organization
		Organization:     service.NewOrganizationService(repo.Organization, log),
		OrganizationType: service.NewOrganizationTypeService(repo.OrganizationType),
		OrgUser:          service.NewOrgUserService(repo.OrgUser, cfg, repo.User), // Cross-repo dependency

		// Terminal & Platform
		Terminal:        service.NewTerminalService(repo.Terminal),
		AppServiceIcon:  service.NewAppServiceIconService(repo.AppServiceIcon),
		AppServiceGroup: service.NewAppServiceIconGroup(repo.AppServiceIconGroup),

		// Content
		PublicFile:   service.NewPublicFileService(repo.PublicFile, cfg),
		Notification: service.NewNotificationService(repo.Notification, cfg),
		News:         service.NewNewsService(repo.News),
		ChatItem:     service.NewChatItemService(repo.ChatItem, log),

		// Logging
		APILog: service.NewAPILogService(repo.APILog),

		// External Integrations
		Verify: service.NewVerifyService(cfg), // XYP, Passport APIs
		Meet:   service.NewMeetService(cfg),   // Video conference API
		Tpay:   service.NewTpayService(cfg),   // Payment API
	}

	// ============================================================
	// STEP 2.5: Initialize Local Auth Services (Redis + Auth)
	// ============================================================
	// Load auth config from environment
	authCfg := localconfig.LoadAuthConfig()

	// Create Redis client for session storage
	redisClient := redis.NewClient(&redis.Options{
		Addr:     authCfg.Redis.Addr(),
		Password: authCfg.Redis.Password,
		DB:       authCfg.Redis.DB,
	})

	// Create Redis session store
	sessionStore := service.NewRedisSessionStore(redisClient, "session:", authCfg.LocalAuth.SessionTTL)
	svc.SessionStore = sessionStore

	// Create Auth service (depends on repo.Auth, sessionStore, and authCfg)
	svc.Auth = service.NewAuthService(repo.Auth, sessionStore, &authCfg.LocalAuth, log)

	// Create Registration service (depends on repo.Auth, repo.User, repo.Registration, svc.Auth)
	svc.Registration = service.NewRegistrationService(
		repo.Auth,
		repo.User,
		repo.Registration,
		svc.Auth,
		&authCfg.LocalAuth,
		log,
	)

	// ============================================================
	// STEP 3: Create permission cache
	// ============================================================
	// Permission cache нь 5 минутын TTL-тэй.
	// Permission шалгахад DB руу дахин дахин очихгүй.
	permCache := auth.NewPermissionCache(permissionSvc, 5*time.Minute)

	// ============================================================
	// STEP 4: Wire up cache invalidators
	// ============================================================
	// Service-ууд permission өөрчлөгдөхөд cache цэвэрлэхэд ашиглана.
	svc.Permission.SetCacheInvalidator(permCache)
	svc.Role.SetCacheInvalidator(permCache)
	svc.UserRole.SetCacheInvalidator(permCache)

	// ============================================================
	// STEP 5: Create final Dependencies struct
	// ============================================================
	return &Dependencies{
		// Core dependencies
		Cfg:       cfg,
		DB:        db,
		Log:       log,
		AuthCache: authCache,

		// SSO client (auth-ийн бүх зүйлийг агуулна)
		SSO: ssoclient.NewSSOClient(cfg, log, authCache),

		// Permission cache (permission шалгахад ашиглана)
		PermCache: permCache,

		// Layer containers
		Repo:    repo,
		Service: svc,
	}
}
