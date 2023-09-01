# DomoticX
This is a personal project whose main objective is to learn about a complete process by developing and deploying a home automation application.

# Table of Contents
1. [Description](#description)
2. [Getting started](#getting-started)
    1. [Dependencies](#dependencies)
    2. [Build](#build)
    3. [Deployment](#deployment)
    4. [Executing program](#executing-program)
3. [Contribute](#contribute)
4. [Authors](#authors)
5. [Version History](#version-history)
6. [License](#license)

## Description
* DomoticX is a comprehensive home automation application developed in Go, aimed at providing a seamless and efficient solution for controlling various smart devices within a household. The project's primary objective is to serve as a learning platform for Go programming language, backend development, and DevOps practices.

## Getting started
```
git clone git@github.com:ltovarm/DomoticX.git
cd DomoticX
make run
```

### Dependencies
* Docker version 20.10.22 or later
* Go 1.20 or later
* Python 3.X.X or later 
* View go.mod and go.sum files for go dependencies

### Build
The make file is in charge of doing all the compiling and raising the services, in this case it has been decided that the docker will be in charge of the compiling, as Go is very light when compiling and this allows us to make the monorepo more isolated and scalable.

### Deployment
All deployment will be implemented in docker, although in the future it will be migrated to AWS or similar.

### Executing program
To run this app it is necessary to run the dockers using the command 
```
make run
```
To stop and delete the containers use the command
```
make stop 
```

## Contribute

## Authors
Contributors names and contact info

Luis Tovar  
[Github](https://github.com/LTovarM)
[Contact](mailto:luistovarmoniz@gmail.com)

## Version History

<!-- * 0.2
    * Various bug fixes and optimizations
    * See [commit change]() or See [release history]() -->
* 0.1
    * Initial Release

## License

This project is licensed under the [Luis Tovar] License - see the LICENSE.md file for details

<!-- ## Acknowledgments

Inspiration, code snippets, etc.
* [awesome-readme](https://github.com/matiassingers/awesome-readme)
* [PurpleBooth](https://gist.github.com/PurpleBooth/109311bb0361f32d87a2)
* [dbader](https://github.com/dbader/readme-template)
* [zenorocha](https://gist.github.com/zenorocha/4526327)
* [fvcproductions](https://gist.github.com/fvcproductions/1bfc2d4aecb01a834b46) -->
