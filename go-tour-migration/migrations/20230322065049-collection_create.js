module.exports = {
  async up(db, client) {
    await db.collection('users').createIndex({ email: 1 }, { unique: true, name: 'email_unique_idx' });
  },

  async down(db, client) {
    await db.collection('users').dropIndex('email_unique_idx');
  }
};
