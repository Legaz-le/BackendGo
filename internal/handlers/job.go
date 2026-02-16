package handlers

import (
	"encoding/json"
	"net/http"
	"strconv"

	"example.com/mod/internal/models"
	"github.com/go-chi/chi/v5"
)

func GetJobs(w http.ResponseWriter, r *http.Request) {
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
	
	var jobs []models.Job
	var err error
	
	if location != "" || minSalary > 0 || maxSalary > 0 {
		jobs, err = models.GetJobsWithFilter(ctx, location, minSalary, maxSalary)
	} else {
		jobs, err = models.GetAllJobs(ctx)
	}
	
	if err != nil {
		http.Error(w, "Internal server error", http.StatusInternalServerError)
		return
	}
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(jobs)
}

func GetJob(w http.ResponseWriter, r *http.Request) {
	id := chi.URLParam(r, "id")
	newId, err := strconv.Atoi(id)
	if err != nil {
		http.Error(w, "Invalid job ID", http.StatusBadRequest)
		return
	}
	job, err := models.GetJobByID(r.Context(), newId)
	if err != nil {
		http.Error(w, "Not found", http.StatusNotFound)
		return
	}
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(job)

}

func PostJob(w http.ResponseWriter, r *http.Request) {
	var job models.Job
	ctx := r.Context()
	if err := json.NewDecoder(r.Body).Decode(&job); err != nil {
		http.Error(w, "Invalid JSON", http.StatusBadRequest)
		return
	}

	if job.Title == "" || job.Description == "" || job.Location == "" {
		http.Error(w, "Missing required fields: title, description, location", http.StatusBadRequest)
		return
	}
	newJob, err := models.CreateJob(ctx, job)
	if err != nil {
		http.Error(w, "Failed to create job", http.StatusInternalServerError)
		return
	}
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(newJob)
}

func PutJob(w http.ResponseWriter, r *http.Request) {
	id := chi.URLParam(r, "id")
	jobID, err := strconv.Atoi(id)
	if err != nil {
		http.Error(w, "Invalid job ID", http.StatusBadRequest)
		return
	}

	var job models.Job
	ctx := r.Context()
	if err := json.NewDecoder(r.Body).Decode(&job); err != nil {
		http.Error(w, "Invalid JSON", http.StatusBadRequest)
		return
	}

	if job.Title == "" || job.Description == "" || job.Location == "" {
		http.Error(w, "Missing required fields: title, description, location", http.StatusBadRequest)
		return
	}

	updatedJob, err := models.UpdateJob(ctx, jobID, job)
	if err != nil {
		http.Error(w, "Job not found or update failed", http.StatusNotFound)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(updatedJob)
}

func DeleteJob(w http.ResponseWriter, r *http.Request) {
	id := chi.URLParam(r, "id")
	jobID, err := strconv.Atoi(id)
	if err != nil {
		http.Error(w, "Invalid job ID", http.StatusBadRequest)
		return
	}

	ctx := r.Context()
	err = models.DeleteJob(ctx, jobID)
	if err != nil {
		http.Error(w, "Job not found", http.StatusNotFound)
		return
	}

	w.WriteHeader(http.StatusNoContent)
}
