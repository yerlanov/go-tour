module.exports = {
  async up(db, client) {
    // expire after 60 seconds
    await db.collection('sessions').createIndex({ expiration: 1 }, { expireAfterSeconds: 60, name: 'expiration_ttl_idx' });

    // expire after 7 days
//    await db.collection('sessions').createIndex({ expiration: 1 }, { expireAfterSeconds: 7 * 24 * 60 * 60, name: 'expiration_ttl_idx' });
    },

  async down(db, client) {
    await db.collection('sessions').dropIndex('expiration_ttl_idx');
  }
};
