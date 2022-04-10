package content

import (
    "arion_shot_api/internal/domain/content"
    "arion_shot_api/internal/platform/auth"
    service "arion_shot_api/internal/services/content"
    "arion_shot_api/internal/utils/cloudinary"
    "context"
    "github.com/cloudinary/cloudinary-go/api/uploader"
    "github.com/pkg/errors"
    "github.com/wgarcia4190/web-router-go/web"
    "mime/multipart"
    "net/http"
    "strings"
)

func UpdateVote(ctx context.Context, writer http.ResponseWriter, request *http.Request) error {
    var voteRequest content.VoteRequest
    if err := web.Decode(request, &voteRequest); err != nil {
        return web.NewRequestError(err, http.StatusBadRequest)
    }

    claims, ok := ctx.Value(auth.Key).(auth.Claims)
    if !ok {
        return web.NewShutdownError("auth claims not in context")
    }

    voteRequest.VoterID = &claims.Subject

    result, err := service.ContentService.UpdateVote(&voteRequest)
    if err != nil {
        return web.NewRequestError(err, http.StatusInternalServerError)
    }

    return web.Respond(ctx, writer, result, http.StatusOK)
}

func CreateContent(ctx context.Context, writer http.ResponseWriter, request *http.Request) error {
    err := request.ParseMultipartForm(32 << 20)
    if err != nil {
        return web.NewRequestError(err, http.StatusInternalServerError)
    }

    claims, ok := ctx.Value(auth.Key).(auth.Claims)
    if !ok {
        return web.NewShutdownError("auth claims not in context")
    }

    challengeID := request.Form.Get("challenge_id")

    if len(strings.TrimSpace(challengeID)) == 0 {
        return web.NewRequestError(errors.New("Challenge ID not provided"), http.StatusBadRequest)
    }

    file, handler, err := request.FormFile("file")
    if err != nil {
        return web.NewRequestError(err, http.StatusInternalServerError)
    }
    defer func(file multipart.File) {
        _ = file.Close()
    }(file)

    clResult, err := cloudinary.CL.Upload.Upload(ctx, file, uploader.UploadParams{
        PublicID: handler.Filename,
    })

    if err != nil {
        return web.NewRequestError(err, http.StatusInternalServerError)
    }

    contentRequest := &content.CreateContentRequest{
        URL:         clResult.SecureURL,
        ChallengeID: challengeID,
        OwnerID:     claims.Subject,
    }

    result, err := service.ContentService.CreateContent(contentRequest)
    if err != nil {
        return web.NewRequestError(err, http.StatusInternalServerError)
    }

    return web.Respond(ctx, writer, result, http.StatusOK)
}
