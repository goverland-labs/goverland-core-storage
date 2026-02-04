# Changelog

All notable changes to this project will be documented in this file.

The format is based on [Keep a Changelog](https://keepachangelog.com/en/1.0.0/), and this project adheres
to [Semantic Versioning](https://semver.org/spec/v2.0.0.html).

## [Unreleased]

## [0.5.3] - 2026-02-04

### Added
- Add support for ENS names in paths with new RPC methods

## [0.5.2] - 2026-02-03

### Fixed
- Save initial token price

## [0.5.1] - 2026-02-03

### Fixed
- Getting delegates for lifetime worker

## [0.5.0] - 2026-02-02

### Changed
- Extend getting DAO delegates based on erc20-votes indexer

### Added
- Processing erc20-votes indexer for delegation events

## [0.4.28] - 2026-02-02

### Added
- Token price when proposal is created

## [0.4.27] - 2025-07-30

### Added
- Check fungible id if proposal is created

### Changed 
- Order top daos only by popularity index

## [0.4.26] - 2025-07-16

### Added
- Send dao to discord for which a token is defined

## [0.4.25] - 2025-07-16

### Added
- Extend delegates summary with votes and proposals counters

## [0.4.24] - 2025-06-30

### Changed
- Top logic

## [0.4.23] - 2025-06-27

### Added 
- Added fungible ids filter to get daos by fungible ids

## [0.4.22] - 2025-06-26

### Changed
- Do not send update dao message if only fungible id is changed
- Set fungible id only if the price is not null

## [0.4.21] - 2025-06-23

### Fixed
- Handling zerion error on token-price worker

## [0.4.20] - 2025-06-19

### Fixed
- Method to recalculate fungible ids

## [0.4.19] - 2025-06-16

### Added
- Verification status for dao
- Method to recalculate fungible ids

## [0.4.18] - 2025-06-07

### Fixed
- Increased fungible worker duration, to avoid Zerion limit reaching  

## [0.4.17] - 2025-05-27

### Changed
- Increase limit for top delegates from 5 to 100
- Fix status response for non existing proposal id
- Extend proposals get by filters endpoint for getting short proposal info 

## [0.4.16] - 2025-05-27

### Added
- Extended fungible info

## [0.4.15] - 2025-04-25

### Changed
- Filter delegates summary by allowlisted dao

## [0.4.14] - 2025-04-23

### Fixed
- Getting delegate votes

## [0.4.13] - 2025-04-16

### Added
- Process delegate skip the vote event

## [0.4.12] - 2025-04-01

### Added
- Populate token price for dao

## [0.4.11] - 2025-04-01

### Added
- Fix proposals cnt calculation

## [0.4.10] - 2025-03-31

### Changed
- Rollback getting voters from case insensitive idx
- Speed up proposals and DAO searches 

## [0.4.9] - 2025-03-26

### Changed
- Rollback getting voters from old idx

## [0.4.8] - 2025-03-25

### Added
- Processing delegates expiration events

### Added
- Retrieve and save token symbol from zerion

### Changed
- Get zerion info for all new daos
- Collect token price 12 times per day
- Add erc20-votes strategy to get fungible info

## [0.4.7] - 2025-03-18

### Added
- Filled fields for core DAO entity
- Added spam field to proposal info
- Added votes subscription

## [0.4.6] - 2025-03-18

### Changed
- Limit for zerion request

## [0.4.5] - 2025-03-17

### Changed
- Update fulltext DAO search by name

## [0.4.4] - 2025-02-03

### Added
- Logic to get token info and token chart from zerion

## [0.4.3] - 2025-01-23

### Added
- Logic to check if verified dao has fungible id in zerion

## [0.4.2] - 2024-12-24

### Added
- Added vote status validation

## [0.4.1] - 2024-12-17

### Added
- Send delegates to core-feed

## [0.4.0] - 2024-12-04

### Added
- Implement getting delegations by address
- Implement getting delegators by address
- Implement getting count of delegations by address

## [0.3.2] - 2024-10-27

### Added
- Calculation of token price

## [0.3.1] - 2024-10-18

### Added
- Register event for "delegate voted"

## [0.3.0] - 2024-10-07

### Added
- Extend DAO object with active proposals ids

## [0.2.7] - 2024-10-07

### Fixed
- Fixed ens names batch creating

## [0.2.6] - 2024-10-02

### Added
- Added vote now logic
- Added filter dao in votes
- Make unique name in ens storage 

## [0.2.5] - 2024-09-22

### Added
- Added delegates total

## [0.2.4] - 2024-09-19

### Added
- Register event for "delegate create proposal”

## [0.2.3] - 2024-09-13

### Added
- Added github actions for building docker image

## [0.2.2] - 2024-09-09

### Added
- Added proxy api for delegates

## [0.2.1] - 2024-09-07

### Fixed
- Change delegates subject

## [0.2.0] - 2024-09-07

### Added
- Delegates consumer

## [0.1.27] - 2024-09-03

### Added
- Expand the calculation of proposal's status taking into account dao succeeded choices

## [0.1.26] - 2024-08-23

### Fixed
- Smart filtering daos by uids and aliases

## [0.1.25] - 2024-08-20

### Changed
- GetByID can search by internal and original DAO identifier

## [0.1.24] - 2024-08-13

### Changed
- Parameter's name for search
- Display no more than 2 proposals from one DAO in top proposal 

## [0.1.23] - 2024-07-30

### Fix
- Sorting proposals list with equal voters count

## [0.1.22] - 2024-07-22

### Added
- Search for votes

## [0.1.21] - 2024-07-06

### Changed
- Top proposals calculation

## [0.1.20] - 2024-07-05

### Changed
- Extend vote response with proposal identifier

## [0.1.19] - 2024-06-27

### Changed
- Speed up votes retrieving

## [0.1.18] - 2024-06-13

### Added
- Recalculate proposals count on delete proposal event

## [0.1.17] - 2024-05-15

### Changed
- Extend DAO recommendations with custom data

## [0.1.16] - 2024-04-27

### Changed
- Fix user votes sorting from "vp desc, created desc" to "created desc"

## [0.1.15] - 2024-04-23
- Update top proposals logic to include not verified daos at the end

## [0.1.14] - 2024-04-21
- Update top proposals logic to include only verified daos and allow 2 proposals from one dao

## [0.1.13] - 2024-04-19

### Added
- Described and implemented protocol for dao recommendations
- Export DAO popularity index

## [0.1.12] - 2024-04-11

### Fixed
- Fixed tests

## [0.1.11] - 2024-04-11

### Fixed
- Fixed dao top iterator
- Fixed linter

## [0.1.10] - 2024-04-11

### Added
- Added top dao cache worker

## [0.1.9] - 2024-04-10

### Added
- Indexes for slow queries
- Daos voted in

## [0.1.8] - 2024-03-22

### Added
- Stats server

## [0.1.7] - 2024-03-20

### Changed
- Skip votes if there is no dao

## [0.1.6] - 2024-03-13

### Added
- Total VP for proposal votes

### Changed
- Ordering for proposal votes

## [0.1.5] - 2024-03-12

### Added
- Added the method to get ens names by addresses

## [0.1.4] - 2024-03-12

### Added
- Store vote after voting

## [0.1.3] - 2024-03-06

### Fixed
- Fixed Dockerfile

## [0.1.2] - 2024-03-02

### Fixed
- Fixed type in protocol directory path

## [0.1.1] - 2024-03-02

### Changed
- Migrated protocol to this repo as different module

## [0.1.0] - 2024-03-01

### Changed
- Changed the path name of the go module
- Updated dependencies versions

### Added
- Added LICENSE information
- Added info for contributing
- Added github issues templates
- Added linter and unit-tests workflows for github actions
- Added badges with link to the license and passed workflows

## [0.0.62] - 2024-02-28

### Added
- Caching the dao top by category for the 5 minutes
- Metrics for producing/publishing messages

## [0.0.61] - 2024-02-22

### Changed
- Reduced proposal top cache interval from 30 to 5 minutes
- Increased number of voters for proposal top query 

## [0.0.60] - 2024-02-13

### Changed
- Storing only unique voters for proposals

## [0.0.59] - 2024-02-09

### Fixed
- Calculating count for proposal top logic

## [0.0.58] - 2024-02-08

### Changed
- Filtering canceled proposals

## [0.0.57] - 2024-02-08

### Changed
- Filtering spam and canceled in counters

## [0.0.56] - 2024-02-08

### Changed
- Send spam field to core-feed
- Send event on proposal deletion event

## [0.0.55] - 2024-02-06

### Added
- Active votes, verified fields to dao 

## [0.0.54] - 2024-02-05

### Changed
- Filters for votes request to allow order by voter's address

## [0.0.53] - 2024-01-29

### Changed
- Filters for votes request to allow the list of proposals and the voter's address

## [0.0.52] - 2023-12-20

### Added
- Spam flag for proposals

### Changed
- Updating timeline for proposals

## [0.0.51] - 2023-12-19

### Changed
- Order Daos by popularity index.

### Added
- Popularity index update

## [0.0.50] - 2023-12-14

### Added
- Implement resolving ens names for voters

## [0.0.49] - 2023-12-06

### Added
- Implement resolving ens names for proposals

## [0.0.48] - 2023-12-04

### Added
- Added voting methods

## [0.0.47] - 2023-11-07

### Changed
- Increase the calculation time for new categories from 1 minute to 1 hour
- Store system categories on DAO updates
- Sorting by state first for searing proposals by query
- Update proposal cnt on every dao update

## [0.0.46] - 2023-11-02

### Added
- Calculating logic for "popular_daos" category

## [0.0.45] - 2023-10-27

### Changed
- Prepare proposal top background caching

## [0.0.44] - 2023-10-26

### Added
- Proposal state calculation

## [0.0.43] - 2023-10-18

### Added
- Proposal count calculation

## [0.0.42] - 2023-10-16

### Fixed
- Decrease max pending count consumers

### Changed
- Actualize DB schema

## [0.0.41] - 2023-10-16

### Fixed
- Dao voters table name

## [0.0.40] - 2023-10-16

### Added
- Dao voters calculation

### Changed
- Top proposals calculation
- Speed up getting votes

## [0.0.39] - 2023-10-04

### Fixed
- Caching proposal top key

## [0.0.38] - 2023-10-04

### Added
- Caching proposal top results

### Changed
- The count's calculation for top proposals

## [0.0.37] - 2023-10-02

### Added
- Caching dao ids provider

## [0.0.36] - 2023-09-18

### Changed
- Increase voting finish window from 1 hour to 6 hours

## [0.0.35] - 2023-09-18

### Added
- Calculating logic for "new_daos" category

## [0.0.34] - 2023-09-13

### Changed
- Logic to select top proposals

## [0.0.33] - 2023-09-12

### Changed
- Mark votes choice field as json.RawMessage due to multiple values

## [0.0.32] - 2023-08-23

### Changed
- Actualize calculating quorum reached event 

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
