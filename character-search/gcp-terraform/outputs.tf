output "mongodb_public_ip" {
  description = "Public IP address of the MongoDB instance"
  value       = google_compute_instance.mongodb_instance.network_interface[0].access_config[0].nat_ip
}
