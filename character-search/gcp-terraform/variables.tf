variable "project_id" {
  description = "GCP Project ID"
  type        = string
}

variable "region" {
  description = "GCP Region"
  type        = string
  default     = "us-central1"
}

variable "zone" {
  description = "GCP Zone"
  type        = string
  default     = "us-central1-a"
}

variable "machine_type" {
  description = "Machine type for the MongoDB instance"
  type        = string
  default     = "e2-medium"
}

variable "credentials_file" {
  description = "Path to the GCP credentials JSON file"
  type        = string
}
