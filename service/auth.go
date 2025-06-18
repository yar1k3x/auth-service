package service

import (
	"AuthService/db"
	"AuthService/model"
	pb "AuthService/proto"
	"AuthService/util"
	"context"
	"log"

	"golang.org/x/crypto/bcrypt"
)

type AuthServiceServer struct {
	pb.UnimplementedAuthServiceServer
}

func (s *AuthServiceServer) Register(ctx context.Context, req *pb.RegisterRequest) (*pb.AuthResponse, error) {
	hashed, _ := bcrypt.GenerateFromPassword([]byte(req.Password), bcrypt.DefaultCost)

	res, err := db.DB.Exec("INSERT INTO users (fio, email, password, role_id) VALUES (?, ?, ?, ?)", req.Fio, req.Email, string(hashed), 1)
	if err != nil {
		log.Println("DB error:", err)
		return &pb.AuthResponse{Error: "Email already used"}, nil
	}
	log.Println("Пользователь", req.Fio, "был создан!")
	id, _ := res.LastInsertId()
	token, _ := util.GenerateJWT(id)
	return &pb.AuthResponse{Token: token}, nil
}

func (s *AuthServiceServer) Login(ctx context.Context, req *pb.LoginRequest) (*pb.AuthResponse, error) {
	user := model.User{}
	row := db.DB.QueryRow("SELECT id, password FROM users WHERE email = ?", req.Email)
	err := row.Scan(&user.ID, &user.Password)
	if err != nil {
		return &pb.AuthResponse{Error: "user not found"}, nil
	}

	if err := bcrypt.CompareHashAndPassword([]byte(user.Password), []byte(req.Password)); err != nil {
		return &pb.AuthResponse{Error: "invalid password"}, nil
	}

	token, _ := util.GenerateJWT(user.ID)
	return &pb.AuthResponse{Token: token}, nil
}

func (s *AuthServiceServer) GetUsers(ctx context.Context, req *pb.GetUsersRequest) (*pb.GetUsersResponse, error) {
	baseQuery := `
		SELECT u.id, u.fio, u.email, ur.name AS role FROM users u
		JOIN user_role ur ON u.role_id = ur.id
		`

	var args []interface{}
	var conditions []string

	if req.UserId != nil {
		conditions = append(conditions, "user_id = ?")
		args = append(args, req.UserId.Value)
	}
	if req.RoleId != nil {
		conditions = append(conditions, "role_id = ?")
		args = append(args, req.RoleId.Value)
	}

	if len(conditions) > 0 {
		baseQuery += " WHERE " + conditions[0]
		for i := 1; i < len(conditions); i++ {
			baseQuery += " AND " + conditions[i]
		}
	}

	rows, err := db.DB.Query(baseQuery, args...)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	var users []*pb.UserTemplate

	for rows.Next() {
		var r pb.UserTemplate
		err := rows.Scan(
			&r.Id,
			&r.Fio,
			&r.Email,
			&r.Role,
		)
		if err != nil {
			return nil, err
		}

		users = append(users, &r)
	}

	if err = rows.Err(); err != nil {
		return nil, err
	}

	return &pb.GetUsersResponse{
		Users: users,
	}, nil
}
