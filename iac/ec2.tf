resource "aws_instance" "instance-1" {
  ami           = local.ami_image
  instance_type = "t4g.nano"

  tags = {
    Name = "Load test blocker EC2"
  }
  key_name                    = aws_key_pair.private_key.key_name
  vpc_security_group_ids      = [aws_security_group.ec2-load-test-blocker-sg.id]
  associate_public_ip_address = true
  availability_zone           = local.availability_zone
  subnet_id                   = aws_subnet.my-load-test-subnet-1a.id
  root_block_device {
    volume_type = "gp2"
    volume_size = 100
  }
  user_data_replace_on_change = true
  user_data                   = templatefile("ec2-boot.tftpl", {})
}

resource "aws_security_group" "ec2-load-test-blocker-sg" {
  name   = "secgrp-load-test-blocker-ec2"
  vpc_id = aws_vpc.my-load-test-vpc.id
  ingress {
    from_port   = 22
    to_port     = 22
    protocol    = "tcp"
    cidr_blocks = ["0.0.0.0/0"]
  }
  ingress {
    from_port   = 8080
    to_port     = 8080
    protocol    = "tcp"
    cidr_blocks = ["0.0.0.0/0"]
  }
  egress {
    from_port   = 0
    to_port     = 0
    protocol    = "-1"
    cidr_blocks = ["0.0.0.0/0"]
  }
}

output "ec2-ssh-connect" {
  value = "'ssh ec2-user@${aws_instance.instance-1.public_ip} -i cert.pem'"
}

output "web-connect" {
  value = "http://${aws_instance.instance-1.public_ip}:8080"
}
