package handlers

import (
	"encoding/json"
	"fmt"
	"go-messages/internal/app"
	"go-messages/internal/models"
	"net/http"

	"github.com/julienschmidt/httprouter"
)

type MessageHandler struct {
	messageService app.MessageService
}

func InitMessageHandlers(m app.MessageService) *MessageHandler {
	return &MessageHandler{
		messageService: m,
	}
}

func (h *MessageHandler) Get(w http.ResponseWriter, r *http.Request, ps httprouter.Params) {
    groupID := ps.ByName("groupID")
    if groupID == "" {
        http.Error(w, "Group ID is required", http.StatusBadRequest)
        return
    }

    req := models.RequestMessages{
        GroupiD: groupID,
        Offset:  0,  // значение по умолчанию
        Count:   10, 
    }

    query := r.URL.Query()
    
    if offsetStr := query.Get("offset"); offsetStr != "" {
        var offset int32
        if _, err := fmt.Sscanf(offsetStr, "%d", &offset); err != nil || offset < 0 {
            http.Error(w, "Invalid offset parameter", http.StatusBadRequest)
            return
        }
        req.Offset = offset
    }

    if countStr := query.Get("count"); countStr != "" {
        var count int32
        if _, err := fmt.Sscanf(countStr, "%d", &count); err != nil || count <= 0 {
            http.Error(w, "Invalid count parameter", http.StatusBadRequest)
            return
        }
        // Ограничиваем максимальное количество
        if count > 100 {
            count = 100
        }
        req.Count = count
    }

    messages, err := h.messageService.Get(r.Context(), req)
    if err != nil {
        http.Error(w, err.Error(), http.StatusInternalServerError)
        return
    }

    // Сериализуем ответ в JSON
    w.Header().Set("Content-Type", "application/json")
    if err := json.NewEncoder(w).Encode(messages); err != nil {
        http.Error(w, "Failed to encode response", http.StatusInternalServerError)
    }
}
