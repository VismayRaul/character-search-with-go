provider "google" {
  project = var.project_id
  region  = var.region
  credentials = file(var.credentials_file)
}

resource "google_compute_network" "vpc_network" {
  name = "mongodb-vpc"
}

resource "google_compute_firewall" "allow_mongodb" {
  name    = "allow-mongodb"
  network = google_compute_network.vpc_network.name

  allow {
    protocol = "tcp"
    ports    = [27017]
  }

  source_ranges = ["0.0.0.0/0"]
}

# Create a Cloud Storage bucket to back up MongoDB data
resource "google_storage_bucket" "mongodb_backup" {
  name          = "mongodb-backups-bucket"
  location      = "US"
  force_destroy = true
}

resource "google_compute_instance" "mongodb_instance" {
  name         = "mongodb-instance"
  machine_type = var.machine_type
  zone         = var.zone

  tags = ["mongodb"]

  boot_disk {
    initialize_params {
      image = "ubuntu-2004-focal-v20241115" # Latest Ubuntu 20.04 LTS image
      size  = 50 # Adjust disk size as needed
    }
  }

  network_interface {
    network = google_compute_network.vpc_network.name
    access_config {}
  }

  metadata_startup_script = <<-EOT
    #!/bin/bash
    apt-get update
    apt-get install -y gnupg
    wget -qO - https://www.mongodb.org/static/pgp/server-6.0.asc | apt-key add -
    echo "deb [ arch=amd64,arm64 ] https://repo.mongodb.org/apt/ubuntu focal/mongodb-org/6.0 multiverse" | tee /etc/apt/sources.list.d/mongodb-org-6.0.list
    apt-get update
    apt-get install -y mongodb-org
    systemctl start mongod
    systemctl enable mongod
  EOT
}

resource "null_resource" "fetch_rick_and_morty_data" {
  depends_on = [google_storage_bucket.mongodb_backup]

  provisioner "local-exec" {
    command = <<-EOT
      curl -o /characters.json https://rickandmortyapi.com/api/character
      if [ $? -ne 0 ]; then
        echo "Failed to fetch Rick and Morty data with curl"
        exit 1
      fi

      gcloud storage cp /characters.json gs://mongodb-backups-bucket/characters.json
      if [ $? -ne 0 ]; then
        echo "Failed to upload data to Google Cloud Storage"
        exit 1
      fi
    EOT
  }
}



output "instance_ip" {
  value = google_compute_instance.mongodb_instance.network_interface[0].access_config[0].nat_ip
}

output "storage_bucket_name" {
  value = google_storage_bucket.mongodb_backup.name
}

