package repos

import (
	"backend/src/db"
	"fmt"
	"time"
)

type User struct{}

func (User) FindByPhoneNumber(phoneNumber string) (*db.User, error) {
	var user db.User
	err := db.DB.Where("phone_number = ?", phoneNumber).First(&user).Error
	if err != nil {
		return nil, err
	}
	return &user, nil
}

func (User) Create(phoneNumber, name string) (*db.User, error) {
	user := &db.User{
		PhoneNumber: phoneNumber,
		Name:        name,
		IsActive:    true,
	}
	if err := db.DB.Create(user).Error; err != nil {
		return nil, err
	}
	return user, nil
}

func (User) GetByID(userID uint) (*db.User, error) {
	var user db.User
	err := db.DB.First(&user, userID).Error
	if err != nil {
		return nil, err
	}
	return &user, nil
}

type OTP struct{}

func (OTP) Create(userID uint, otpCode string, expiresAt int64) (*db.UserOtp, error) {
	otp := &db.UserOtp{
		UserID:     userID,
		OtpCode:    otpCode,
		ExpiresAt:  expiresAt,
		IsVerified: false,
	}
	if err := db.DB.Create(otp).Error; err != nil {
		return nil, err
	}
	return otp, nil
}

func (OTP) FindByID(otpID uint) (*db.UserOtp, error) {
	var otp db.UserOtp
	err := db.DB.First(&otp, otpID).Error
	if err != nil {
		return nil, err
	}
	return &otp, nil
}

func (OTP) Verify(otpID uint, otpCode string) error {
	var otp db.UserOtp
	if err := db.DB.First(&otp, otpID).Error; err != nil {
		return err
	}

	if otp.IsVerified {
		return fmt.Errorf("OTP already verified")
	}

	if otp.ExpiresAt < time.Now().Unix() {
		return fmt.Errorf("OTP expired")
	}

	if otp.OtpCode != otpCode {
		return fmt.Errorf("invalid OTP")
	}

	otp.IsVerified = true
	return db.DB.Save(&otp).Error
}

func (OTP) CountRecentByPhoneNumber(phoneNumber string, sinceMinutes int) (int64, error) {
	var count int64
	since := time.Now().Add(-time.Duration(sinceMinutes) * time.Minute).Unix()

	err := db.DB.Model(&db.UserOtp{}).
		Joins("JOIN users ON users.id = user_otps.user_id").
		Where("users.phone_number = ? AND user_otps.created_at > ?", phoneNumber, since).
		Count(&count).Error

	return count, err
}

type SessionRepo struct{}

func (SessionRepo) Create(userID uint, sessionID string, expiresAt int64, ipAddress, userAgent string) (*db.Session, error) {
	session := &db.Session{
		SessionID: sessionID,
		UserID:    userID,
		ExpiresAt: expiresAt,
		IPAddress: ipAddress,
		UserAgent: userAgent,
	}
	if err := db.DB.Create(session).Error; err != nil {
		return nil, err
	}
	return session, nil
}

func (SessionRepo) FindBySessionID(sessionID string) (*db.Session, error) {
	var session db.Session
	err := db.DB.Preload("User").Where("session_id = ? AND expires_at > ?", sessionID, time.Now().Unix()).First(&session).Error
	if err != nil {
		return nil, err
	}
	return &session, nil
}

func (SessionRepo) Delete(sessionID string) error {
	return db.DB.Where("session_id = ?", sessionID).Delete(&db.Session{}).Error
}

func (SessionRepo) DeleteExpired() error {
	return db.DB.Where("expires_at < ?", time.Now().Unix()).Delete(&db.Session{}).Error
}

func (SessionRepo) DeleteByUserID(userID uint) error {
	return db.DB.Where("user_id = ?", userID).Delete(&db.Session{}).Error
}

type UserWithSession struct {
	ID             uint
	PhoneNumber    string
	Name           string
	Email          string
	SavingsBalance int
	SharesBalance  int
	IsActive       bool
}
