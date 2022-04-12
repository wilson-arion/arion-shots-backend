package user

import (
    "arion_shot_api/domain/user"
    "arion_shot_api/platform/auth"
    sqlutils "arion_shot_api/utils/sql"
    "github.com/pkg/errors"
    "time"
)

const (
    queryFindByEmailAndPassword = `
        SELECT
	        user_id,
	        first_name,
	        last_name,
	        email,
            user_role,
	        date_created,
	        date_updated
        FROM
	        users
        WHERE
	        email = ? and pass = MD5(?);
    `

    queryEmailExist = `
        SELECT
            IF(COUNT(*) > 0, 'true', 'false') AS value
        FROM
            users
        WHERE
            email = ?
    `

    queryCreateUser = `
        INSERT INTO users
            (user_id, first_name, last_name, email, pass, user_role, date_created, date_updated)
        VALUES
            (0, ?, ?, ?, MD5(?), ?, now(), now());
    `
)

var (
    UserRepository userRepositoryInterface = &userRepository{}
)

type userRepository struct{}

type userRepositoryInterface interface {
    FindByEmailAndPassword(request user.LoginRequest) (*user.User, error)
    CreateUser(request user.RegisterRequest) (*user.User, error)
    EmailExist(email string) bool
    Authenticate(userId string, role string, now time.Time) (auth.Claims, error)
}

func (repository *userRepository) FindByEmailAndPassword(request user.LoginRequest) (*user.User, error) {
    stmt, err := sqlutils.CreateStmt(queryFindByEmailAndPassword) //nolint:sqlclosecheck

    if err != nil {
        return nil, err
    }
    defer sqlutils.CloseStmt(stmt)

    u := &user.User{}

    result := stmt.QueryRow(request.Email, request.Password)
    if err := result.Scan(
        &u.ID,
        &u.FirstName,
        &u.LastName,
        &u.Email,
        &u.Role,
        &u.DateCreated,
        &u.DateUpdated,
    ); err != nil {
        return nil, sqlutils.ParseError(err)
    }

    println(request.GetEncryptedPassword())

    return u, nil
}

func (repository *userRepository) CreateUser(request user.RegisterRequest) (*user.User, error) {
    stmt, err := sqlutils.CreateStmt(queryCreateUser) //nolint:sqlclosecheck

    if err != nil {
        return nil, err
    }
    defer sqlutils.CloseStmt(stmt)

    insertResult, insertErr := stmt.Exec(request.FirstName, request.LastName, request.Email, request.Password, request.Role)
    if insertErr != nil {
        return nil, sqlutils.ParseError(insertErr)
    }

    if insertResult != nil {
        response, err := repository.FindByEmailAndPassword(user.LoginRequest{
            Email:    request.Email,
            Password: request.Password,
        })
        if err != nil {
            return nil, err
        }

        return response, nil
    }

    return nil, errors.New("Error creating new user")
}

func (repository *userRepository) EmailExist(email string) bool {
    stmt, err := sqlutils.CreateStmt(queryEmailExist) //nolint:sqlclosecheck

    if err != nil {
        return false
    }
    defer sqlutils.CloseStmt(stmt)

    var tmp bool

    result := stmt.QueryRow(email)
    if err := result.Scan(&tmp); err != nil {
        return false
    }

    return tmp
}

func (repository *userRepository) Authenticate(userId string, role string, now time.Time) (auth.Claims, error) {
    // If we are this far the request is valid. Create some claims from the user
    // and generate their token.
    claims := auth.NewClaims(userId, []string{role}, now, time.Hour)
    return claims, nil
}
