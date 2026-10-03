const fs = require('fs');

print("🚀 Checking MongoDB for existing data...");

const seedDb = db.getSiblingDB("magic-stream-movies");

// Check if already seeded
const hasData =
  seedDb.movies.estimatedDocumentCount() > 0 ||
  seedDb.genres.estimatedDocumentCount() > 0 ||
  seedDb.users.estimatedDocumentCount() > 0 ||
  seedDb.rankings.estimatedDocumentCount() > 0;

if (hasData) {
  print("⚠️ Data already exists — skipping seed.");
  quit();
}

print("📌 Seeding initial database data...");

seedDb.users.insertMany(EJSON.parse(fs.readFileSync('/seed-data/users.json', 'utf8')));
seedDb.movies.insertMany(EJSON.parse(fs.readFileSync('/seed-data/movies.json', 'utf8')));
seedDb.genres.insertMany(EJSON.parse(fs.readFileSync('/seed-data/genres.json', 'utf8')));
seedDb.rankings.insertMany(EJSON.parse(fs.readFileSync('/seed-data/rankings.json', 'utf8')));

print("🎉 Seeding complete. Database is ready!");