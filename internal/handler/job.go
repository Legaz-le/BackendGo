package handler

import (
	"encoding/json"
	"net/http"
	"strconv"

	"example.com/mod/internal/domain"
	"example.com/mod/internal/service"
	"github.com/go-chi/chi/v5"
)

type JobHandler struct {
	service *service.JobService
}

func NewJobHandler(service *service.JobService) *JobHandler {
	return &JobHandler{service: service}
}

func (h *JobHandler) GetJobs(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()
	location := r.URL.Query().Get("location")
	minSalaryStr := r.URL.Query().Get("salary_min")
	maxSalaryStr := r.URL.Query().Get("salary_max")

	minSalary := 0
	maxSalary := 0

	if minSalaryStr != "" {
		if val, err := strconv.Atoi(minSalaryStr); err == nil {
			minSalary = val
		}
	}

	if maxSalaryStr != "" {
		if val, err := strconv.Atoi(maxSalaryStr); err == nil {
			maxSalary = val
		}
	}

	var jobs []domain.Job
	var err error

	if location != "" || minSalary > 0 || maxSalary > 0 {
		jobs, err = h.service.GetWithFilter(ctx, location, minSalary, maxSalary)
	} else {
		jobs, err = h.service.GetAll(ctx)
	}

	if err != nil {
		http.Error(w, "Internal server error", http.StatusInternalServerError)
		return
	}
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(jobs)
}

func (h *JobHandler) GetJob(w http.ResponseWriter, r *http.Request) {
	id := chi.URLParam(r, "id")
	newId, err := strconv.Atoi(id)
	if err != nil {
		http.Error(w, "Invalid job ID", http.StatusBadRequest)
		return
	}
	job, err := h.service.GetByID(r.Context(), newId)
	if err != nil {
		http.Error(w, "Not found", http.StatusNotFound)
		return
	}
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(job)

}

func (h *JobHandler) PostJob(w http.ResponseWriter, r *http.Request) {
	var job domain.Job
	ctx := r.Context()
	if err := json.NewDecoder(r.Body).Decode(&job); err != nil {
		http.Error(w, "Invalid JSON", http.StatusBadRequest)
		return
	}

	if job.Title == "" || job.Description == "" || job.Location == "" {
		http.Error(w, "Missing required fields: title, description, location", http.StatusBadRequest)
		return
	}
	newJob, err := h.service.Create(ctx, job)
	if err != nil {
		http.Error(w, "Failed to create job", http.StatusInternalServerError)
		return
	}
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(newJob)
}

func (h *JobHandler) PutJob(w http.ResponseWriter, r *http.Request) {
	id := chi.URLParam(r, "id")
	jobID, err := strconv.Atoi(id)
	if err != nil {
		http.Error(w, "Invalid job ID", http.StatusBadRequest)
		return
	}

	var job domain.Job
	ctx := r.Context()
	if err := json.NewDecoder(r.Body).Decode(&job); err != nil {
		http.Error(w, "Invalid JSON", http.StatusBadRequest)
		return
	}

	if job.Title == "" || job.Description == "" || job.Location == "" {
		http.Error(w, "Missing required fields: title, description, location", http.StatusBadRequest)
		return
	}

	updatedJob, err := h.service.Update(ctx, jobID, job)
	if err != nil {
		http.Error(w, "Job not found or update failed", http.StatusNotFound)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(updatedJob)
}

func (h *JobHandler) DeleteJob(w http.ResponseWriter, r *http.Request) {
	id := chi.URLParam(r, "id")
	jobID, err := strconv.Atoi(id)
	if err != nil {
		http.Error(w, "Invalid job ID", http.StatusBadRequest)
		return
	}

	ctx := r.Context()
	err = h.service.Delete(ctx, jobID)
	if err != nil {
		http.Error(w, "Job not found", http.StatusNotFound)
		return
	}

	w.WriteHeader(http.StatusNoContent)
}

func HealthCheck(w http.ResponseWriter, r *http.Request) {
	w.WriteHeader(http.StatusOK)
}
