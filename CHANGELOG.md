# Changelog

All notable changes to this project will be documented in this file.

The format is based on [Keep a Changelog](https://keepachangelog.com/en/1.0.0/), and this project adheres
to [Semantic Versioning](https://semver.org/spec/v2.0.0.html).

## [Unreleased]

## [0.0.30] - 2023-08-23

### Changed
- Extend proposal grpc response with timeline field 

### Added
- Store proposal timeline to the database

## [0.0.29] - 2023-08-14

### Added
- Register proposal ends soon event

## [0.0.28] - 2023-07-18

### Changed
- Extend vote model

## [0.0.27] - 2023-07-15

### Fixed
- Updated platform-events dependency to v0.0.20
- Fixed handling dao activity since event for new dao
- Fixed duplicate errors on getting an internal id for dao
- Fixed payload from proposal handling for core.dao.check.activity_since event

### Changed
- Disabled not found errors in gorm logger

## [0.0.26] - 2023-07-15

### Fixed
- Updated platform-events dependency to v0.0.17

## [0.0.25] - 2023-07-14

### Fixed
- Updated platform-events dependency to v0.0.16

## [0.0.24] - 2023-07-14

### Fixed
- Fixed activity since processing

## [0.0.23] - 2023-07-14

### Fixed
- Fixed issues after rebasing

## [0.0.22] - 2023-07-14

### Added
- Activity since to the DAO model

## [0.0.21] - 2023-07-14

### Fixed
- Updated platform-events dependency to v0.0.14

## [0.0.20] - 2023-07-14

### Changed
- Sorting dao list by followers 
- Sorting proposals list by voters

## [0.0.19] - 2023-07-13

### Fixed
- Fixed events for vote starts soon, vote started, vote ended 
- Fixed filter DAOs by name

## [0.0.18] - 2023-07-12

### Fixed
- Updated platform-events dependency to v0.0.13

## [0.0.17] - 2023-07-12

### Fixed
- Fixed comparing proposals in consumer for updating

## [0.0.16] - 2023-07-12

### Fixed
- Fixed missed fields in DAO and Strategy params
- Updated core-api protocol version to v0.0.8

## [0.0.15] - 2023-07-11

### Fixed
- Fixed missed fields in DAO and Strategy objects
- Updated platform-events dependency to v0.0.12

## [0.0.14] - 2023-07-11

### Fixed
- Fixed tests for update dao if needed
- Updated platform-events dependency to v0.0.11

## [0.0.13] - 2023-07-11

### Added
- Order by votes proposal filter

## [0.0.12] - 2023-07-07

### Added
- Filtering proposals by title

## [0.0.11] - 2023-07-06

### Fixed
- Fixed error checks in getting dao id 

## [0.0.10] - 2023-07-06

### Fixed
- Fixed creation dao id service

## [0.0.9] - 2023-07-06

### Added
- Added dao id provider for generating UUID ids

## [0.0.8] - 2023-07-06

### Fixed
- Fixed dockerfile ans infra files structure

## [0.0.7] - 2023-07-04

### Changed
- DAO id generation

## [0.0.6] - 2023-06-29

### Added
- Filter dao by ids

## [0.0.5] - 2023-06-14

### Added
- Dao gRPC server 
- Proposal gRPC server 
- Vote gRPC server 

## [0.0.4] - 2023-05-23

### Added
- Configure debugging DB queries by env

### Changed
- Using group name in consumers

## [0.0.3] - 2023-05-18

### Added
- Add vote handling
- Describe basic internal events

## [0.0.2] - 2023-04-25

### Added
- Dockerfile

## [0.0.1] - 2023-04-25

### Added
- Basic schema for storing daos
- Basic schema for storing proposals
