package handler

import (
	"net/http"

	"pet-platform/backend/internal/domain"
	"pet-platform/backend/internal/service"
)

type AuthHandler struct {
	authService service.AuthService
}

func NewAuthHandler(authService service.AuthService) *AuthHandler {
	return &AuthHandler{authService: authService}
}

func (h *AuthHandler) Register(w http.ResponseWriter, r *http.Request) {
	var req domain.RegisterRequest
	if err := readJSON(r, &req); err != nil {
		writeError(w, http.StatusBadRequest, "BAD_REQUEST", "Invalid request body")
		return
	}

	if req.Email == "" || req.Password == "" || req.FirstName == "" {
		writeError(w, http.StatusBadRequest, "VALIDATION_FAILED", "Email, password, and first name are required")
		return
	}

	res, err := h.authService.Register(r.Context(), &req)
	if err != nil {
		writeError(w, http.StatusConflict, "REGISTRATION_FAILED", err.Error())
		return
	}

	writeJSON(w, http.StatusCreated, res)
}

func (h *AuthHandler) Login(w http.ResponseWriter, r *http.Request) {
	var req domain.LoginRequest
	if err := readJSON(r, &req); err != nil {
		writeError(w, http.StatusBadRequest, "BAD_REQUEST", "Invalid request body")
		return
	}

	res, err := h.authService.Login(r.Context(), &req)
	if err != nil {
		writeError(w, http.StatusUnauthorized, "UNAUTHORIZED", err.Error())
		return
	}

	writeJSON(w, http.StatusOK, res)
}

func (h *AuthHandler) Refresh(w http.ResponseWriter, r *http.Request) {
	var body struct {
		RefreshToken string `json:"refreshToken"`
	}
	if err := readJSON(r, &body); err != nil || body.RefreshToken == "" {
		writeError(w, http.StatusBadRequest, "BAD_REQUEST", "refreshToken is required")
		return
	}

	tokens, err := h.authService.RefreshToken(r.Context(), body.RefreshToken)
	if err != nil {
		writeError(w, http.StatusUnauthorized, "UNAUTHORIZED", err.Error())
		return
	}

	writeJSON(w, http.StatusOK, tokens)
}

func (h *AuthHandler) Me(w http.ResponseWriter, r *http.Request) {
	userID, ok := getUserID(r)
	if !ok {
		writeError(w, http.StatusUnauthorized, "UNAUTHORIZED", "User not authenticated")
		return
	}

	user, err := h.authService.GetProfile(r.Context(), userID)
	if err != nil || user == nil {
		writeError(w, http.StatusNotFound, "NOT_FOUND", "User profile not found")
		return
	}

	writeJSON(w, http.StatusOK, user)
}

func (h *AuthHandler) RequestOTP(w http.ResponseWriter, r *http.Request) {
	var body struct {
		Target  string `json:"target"`
		Purpose string `json:"purpose"`
	}
	if err := readJSON(r, &body); err != nil || body.Target == "" {
		writeError(w, http.StatusBadRequest, "BAD_REQUEST", "target is required")
		return
	}

	_ = h.authService.RequestOTP(r.Context(), body.Target, body.Purpose)
	writeSuccess(w, http.StatusOK, "OTP sent if account exists")
}

func (h *AuthHandler) VerifyOTP(w http.ResponseWriter, r *http.Request) {
	var body struct {
		Target  string `json:"target"`
		OTP     string `json:"otp"`
		Purpose string `json:"purpose"`
	}
	if err := readJSON(r, &body); err != nil || body.Target == "" || body.OTP == "" {
		writeError(w, http.StatusBadRequest, "BAD_REQUEST", "target and otp are required")
		return
	}

	res, err := h.authService.VerifyOTP(r.Context(), body.Target, body.OTP, body.Purpose)
	if err != nil {
		writeError(w, http.StatusUnauthorized, "INVALID_OTP", err.Error())
		return
	}

	writeJSON(w, http.StatusOK, res)
}
