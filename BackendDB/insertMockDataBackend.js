
const faker = require('faker');
const admin = require('firebase-admin');
const { getFirestore } = require('firebase-admin/firestore');

const serviceAccount = require('./unischedule-5ee93-firebase-adminsdk-ci6y9-6cde344deb.json');


const db = getFirestore(admin.initializeApp({
    credential: admin.credential.cert(serviceAccount)
  }),'unischedule-backend');




  const NUM_USERS = 24;
  const NUM_GROUPS = 12;
  const EVENTS_PER_USER = 5;
  const EVENTS_PER_GROUP = 5;
  
  const createUser = async () => {
    const userRef = db.collection('Users').doc();
    await userRef.set({
      name: faker.name.findName(),
      profilePicture: 'https://picsum.photos/350',
      friends: [],
      groups: [],
      events: []
    });
    return userRef.id;
  };
  
  const createGroup = async (possibleMembers) => {
    const selectedMembers = faker.random.arrayElements(possibleMembers, faker.random.number({min: 6, max: 12}));
    const groupRef = db.collection('Groups').doc();
    const emojis = ['😊', '👍', '🌟', '🎉', '🔥', '💻', '🎨', '📚', '⚽', '🎵']; 
    const randomEmoji = faker.random.arrayElement(emojis); 

    await groupRef.set({
      name: faker.company.companyName(),
      groupPicture: 'https://picsum.photos/1080',
      color: faker.internet.color(),
      icon: randomEmoji,
      members: selectedMembers,
      events: []
    });
    selectedMembers.forEach(async (userId) => {
      await db.collection('Users').doc(userId).update({
        groups: admin.firestore.FieldValue.arrayUnion(groupRef.id)
      });
    });
  
    return groupRef.id;
  };
  
  const createEvent = async () => {
    const startDate = faker.date.future();
    const endDate = new Date(startDate.getTime() + faker.random.number({min: 30, max: 180}) * 60000); 
    const eventRef = db.collection('Events').doc();
    await eventRef.set({
      name: faker.lorem.words(),
      description: faker.lorem.sentence(),
      startDate,
      endDate,
      color: faker.internet.color(),
      labels: [faker.lorem.word(), faker.lorem.word()],
      reminder: faker.random.number({min: 5, max: 120})
    });
    return eventRef.id;
  };
  
  const assignEventsAndFriendsToUsers = async (userIds) => {
    for (const userId of userIds) {
      const eventIds = [];
      for (let i = 0; i < EVENTS_PER_USER; i++) {
        const eventId = await createEvent();
        eventIds.push(eventId);
      }
  
      const friendIds = faker.random.arrayElements(userIds.filter(id => id !== userId), faker.random.number({min: 10, max: 12}));
  
      await db.collection('Users').doc(userId).update({
        events: eventIds,
        friends: friendIds
      });
    }
  };
  
  const assignEventsToGroups = async (groupIds) => {
    for (const groupId of groupIds) {
      const eventIds = [];
      for (let i = 0; i < EVENTS_PER_GROUP; i++) {
        const eventId = await createEvent();
        eventIds.push(eventId);
      }
  
      await db.collection('Groups').doc(groupId).update({
        events: eventIds
      });
    }
  };
  
  const populateFirestore = async () => {
    try {
      const userPromises = Array.from({length: NUM_USERS}).map(createUser);
      const userIds = await Promise.all(userPromises);
  
      const groupPromises = [];
      for (let i = 0; i < NUM_GROUPS; i++) {
        groupPromises.push(createGroup(userIds));
      }
      
      const groupIds = await Promise.all(groupPromises);
  
      const eventPromises = [];
      for (let i = 0; i < NUM_USERS * EVENTS_PER_USER + NUM_GROUPS * EVENTS_PER_GROUP; i++) {
        eventPromises.push(createEvent());
      }
      await Promise.all(eventPromises);
  
      await assignEventsAndFriendsToUsers(userIds);
      await assignEventsToGroups(groupIds);
  
      console.log('Firestore successfully populated');
    } catch (error) {
      console.error('Error populating Firestore:', error);
    }
  };
  
  populateFirestore();