package app

import (
	"testing"

	"github.com/ca17/teamsacs/common"
	"github.com/ca17/teamsacs/models"
	"gorm.io/driver/sqlite"
	"gorm.io/gorm"
	"gorm.io/gorm/schema"
)

func newTestApp(t *testing.T) *Application {
	t.Helper()
	db, err := gorm.Open(sqlite.Open("file::memory:?cache=shared"), &gorm.Config{
		NamingStrategy: schema.NamingStrategy{SingularTable: true},
	})
	if err != nil {
		t.Fatalf("open sqlite: %v", err)
	}
	if err := db.AutoMigrate(&models.SysOpr{}, &models.SysConfig{}, &models.NetNode{}); err != nil {
		t.Fatalf("migrate: %v", err)
	}
	return &Application{gormDB: db}
}

func TestEnsureBootstrapAdminCreatesOnce(t *testing.T) {
	t.Setenv("TEAMSACS_BOOTSTRAP_ADMIN_USERNAME", "admin")
	t.Setenv("TEAMSACS_BOOTSTRAP_ADMIN_PASSWORD", "first-secret")
	a := newTestApp(t)

	if err := a.ensureBootstrapAdmin(); err != nil {
		t.Fatalf("first bootstrap: %v", err)
	}

	var count int64
	a.gormDB.Model(&models.SysOpr{}).Where("username = ?", "admin").Count(&count)
	if count != 1 {
		t.Fatalf("count after create = %d, want 1", count)
	}

	var user models.SysOpr
	if err := a.gormDB.Where("username = ?", "admin").First(&user).Error; err != nil {
		t.Fatalf("load user: %v", err)
	}
	wantHash := common.Sha256HashWithSalt("first-secret", common.SecretSalt)
	if user.Password != wantHash {
		t.Fatalf("password hash mismatch after create")
	}

	// Simulate production upgrade / restart with a different bootstrap password.
	t.Setenv("TEAMSACS_BOOTSTRAP_ADMIN_PASSWORD", "should-not-overwrite")
	if err := a.ensureBootstrapAdmin(); err != nil {
		t.Fatalf("second bootstrap: %v", err)
	}
	a.gormDB.Model(&models.SysOpr{}).Where("username = ?", "admin").Count(&count)
	if count != 1 {
		t.Fatalf("count after second bootstrap = %d, want 1", count)
	}
	if err := a.gormDB.Where("username = ?", "admin").First(&user).Error; err != nil {
		t.Fatalf("reload user: %v", err)
	}
	if user.Password != wantHash {
		t.Fatal("existing production password was overwritten")
	}
}

func TestEnsureBootstrapAdminSkipsCustomAdmin(t *testing.T) {
	t.Setenv("TEAMSACS_BOOTSTRAP_ADMIN_USERNAME", "admin")
	t.Setenv("TEAMSACS_BOOTSTRAP_ADMIN_PASSWORD", "default-should-not-apply")
	a := newTestApp(t)

	customHash := common.Sha256HashWithSalt("already-rotated", common.SecretSalt)
	if err := a.gormDB.Create(&models.SysOpr{
		ID:       42,
		Username: "admin",
		Password: customHash,
		Level:    "super",
		Status:   "enabled",
	}).Error; err != nil {
		t.Fatalf("seed existing admin: %v", err)
	}

	if err := a.ensureBootstrapAdmin(); err != nil {
		t.Fatalf("bootstrap: %v", err)
	}
	var user models.SysOpr
	if err := a.gormDB.Where("username = ?", "admin").First(&user).Error; err != nil {
		t.Fatalf("load user: %v", err)
	}
	if user.Password != customHash {
		t.Fatal("custom existing password was changed")
	}
}
