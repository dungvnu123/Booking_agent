package notification

// TODO: implement the repository with parameterized pgx SQL after the real schema
// has been compared with entities and manual SQL in the original Java project.
// Use pgx.Tx for transactional calls; do not create an independent transaction
// inside every repository method. No queries or schema changes run from this file.
