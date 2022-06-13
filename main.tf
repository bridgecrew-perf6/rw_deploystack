variable "project_id" {
  type = string
}

variable "project_number" {
  type = string
}

variable "region" {
  type = string
}

variable "zone" {
  type = string
}

variable "basename" {
    type = string
}

variable "nodes" {
    type = number
}

resource "google_compute_instance" "vm_instance" {
    name = "rwtest"
    machine_type = "f1-micro"
    zone = "us-central1-a"
    boot_disk {
    initialize_params {
        image = "centos-cloud/centos-7"
    }
}
    network_interface {
        subnetwork = "default"
    }
}