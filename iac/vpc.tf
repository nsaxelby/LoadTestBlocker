resource "aws_vpc" "my-load-test-vpc" {
  cidr_block = "10.230.0.0/16"
  tags = {
    Name = "my-load-test-vpc"
  }
}

resource "aws_subnet" "my-load-test-subnet-1a" {
  availability_zone = local.availability_zone
  cidr_block        = "10.230.0.0/24"
  vpc_id            = aws_vpc.my-load-test-vpc.id
  tags = {
    Name = "my-load-test-subnet-1a"
  }
}

resource "aws_internet_gateway" "my-load-test-igw" {
  vpc_id = aws_vpc.my-load-test-vpc.id
  tags = {
    Name = "my-load-test-intnet-gateway"
  }
}

resource "aws_route_table" "my-load-test-subnet-rt" {
  vpc_id = aws_vpc.my-load-test-vpc.id
  tags = {
    Name = "my-load-test-route-table"
  }

  route {
    cidr_block = "0.0.0.0/0"
    gateway_id = aws_internet_gateway.my-load-test-igw.id
  }
}

resource "aws_route_table_association" "route-table-assoc-sub-1a" {
  subnet_id      = aws_subnet.my-load-test-subnet-1a.id
  route_table_id = aws_route_table.my-load-test-subnet-rt.id
}
