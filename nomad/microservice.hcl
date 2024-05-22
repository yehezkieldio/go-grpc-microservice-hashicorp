job "microservices" {
    datacenters = ["dc1"]

    group "inventory" {
         network {
            port "grpc" {
                to = 50051
            }
        }
        task "inventory" {
            driver = "docker"
            service {
                name = "inventory"
                provider = "nomad"
                port = "grpc"
            }
            config {
                image = "ghcr.io/yehezkieldio/go-grpc-microservice-hashicorp/ggmh-inventory"
                command = "/app/inventory"
                ports = ["grpc"]
            }
            resources {
                cpu    = 10
                memory = 50
            }
        }
    }
}