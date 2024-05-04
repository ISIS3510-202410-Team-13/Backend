const express = require('express');
const admin = require('firebase-admin');
const { getFirestore } = require('firebase-admin/firestore');
const { BigQuery } = require('@google-cloud/bigquery');
const serviceAccount = require('../unischedule-5ee93-firebase-adminsdk-ci6y9-6cde344deb.json');

const firestore = getFirestore(admin.initializeApp({
  credential: admin.credential.cert(serviceAccount)
}),'unischedule-backend');
const bigqueryClient = new BigQuery({
  projectId: 'unischedule-5ee93',
  keyFilename: './unischedule-5ee93-ac56b2db69b3.json'
});

const app = express();
app.use(express.json());

app.get('/user/:id/groups', async (req, res) => {
  const { id } = req.params;
  const userRef = await firestore.collection('Users').doc(id).get();

  if (!userRef.exists) {
    return res.status(404).send('User not found');
  }

  const userData = userRef.data();
  const groups = await Promise.all(userData.groups.map(async (groupId) => {
    const groupRef = await firestore.collection('Groups').doc(groupId).get();
    if (!groupRef.exists) return null;

    const groupData = groupRef.data();
    const memberProfiles = await Promise.all(groupData.members.slice(0, 3).map(async (memberId) => {
      const memberRef = await firestore.collection('Users').doc(memberId).get();
      return memberRef.exists ? memberRef.data().profilePicture : null;
    }));

    return {
      id: groupId,
      ...groupData,
      profilePictures: memberProfiles.filter(pic => pic !== null),
      memberCount: groupData.members.length
    };
  }));

  res.status(200).json(groups.filter(group => group !== null));
});

app.get('/user/:id/friends', async (req, res) => {
  const { id } = req.params;
  const userRef = await firestore.collection('Users').doc(id).get();

  if (!userRef.exists) {
    return res.status(404).send('User not found');
  }

  const userData = userRef.data();
  const friends = await Promise.all(userData.friends.map(async (friendId) => {
    const friendRef = await firestore.collection('Users').doc(friendId).get();
    if (!friendRef.exists) return null;

    const friendData = friendRef.data();
    return {
      id: friendId,
      name: friendData.name,
      profilePicture: friendData.profilePicture
    };
  }));

  res.status(200).json(friends.filter(friend => friend !== null));
});

app.get('/user/:id/events', async (req, res) => {
  const { id } = req.params;
  const userRef = await firestore.collection('Users').doc(id).get();

  if (!userRef.exists) {
    return res.status(404).send('User not found');
  }

  const userData = userRef.data();
  const events = await Promise.all(userData.events.map(async (eventId) => {
    const eventRef = await firestore.collection('Events').doc(eventId).get();
    return eventRef.exists ? { id: eventId, ...eventRef.data() } : null;
  }));

  res.status(200).json(events.filter(event => event !== null));
});

app.get('/group/:id/members', async (req, res) => {
    const { id } = req.params;
    const groupRef = await firestore.collection('Groups').doc(id).get();
  
    if (!groupRef.exists) {
      return res.status(404).send('Group not found');
    }
  
    const groupData = groupRef.data();
    const members = await Promise.all(groupData.members.map(async (memberId) => {
      const memberRef = await firestore.collection('Users').doc(memberId).get();
      if (!memberRef.exists) return null;
  
      const memberData = memberRef.data();
      return {
        id: memberId,
        profilePicture: memberData.profilePicture,
        name: memberData.name,
        events: memberData.events
      };
    }));
  
    res.status(200).json(members.filter(member => member !== null));
  });
  
app.get('/group/:id/events', async (req, res) => {
  const { id } = req.params;
  const groupRef = await firestore.collection('Groups').doc(id).get();

  if (!groupRef.exists) {
    return res.status(404).send('Group not found');
  }

  const groupData = groupRef.data();
  const events = await Promise.all(groupData.events.map(async (eventId) => {
    const eventRef = await firestore.collection('Events').doc(eventId).get();
    return eventRef.exists ? { id: eventId, ...eventRef.data() } : null;
  }));

  res.status(200).json(events.filter(event => event !== null));
});

app.post('/user/:id/events', async (req, res) => {
  const { id } = req.params;
  const eventData = req.body;
  const userRef = firestore.collection('Users').doc(id);
  const userSnapshot = await userRef.get();
  if (!userSnapshot.exists) {
      return res.status(404).send('User not found');
  }
  const eventRef = await firestore.collection('Events').add({
      ...eventData,
  });

  await userRef.update({
      events: admin.firestore.FieldValue.arrayUnion(eventRef.id)
  });

  res.status(201).json({
      message: 'Event added successfully',
      eventId: eventRef.id
  });
});

app.post('/analytics/:id/button-tap', async (req, res) => {
  const { id } = req.params;
  const { buttonName, timestamp, userType } = req.body;
  const datasetId = 'unischedule_analytics';
  const tableId = 'button_taps';
  const rows = [{ buttonName, timestamp, userType, userId: id }];
  try {
    await bigqueryClient.dataset(datasetId).table(tableId).insert(rows);
    res.status(201).send('Button tap recorded');
  }
  catch (error) {
    console.error('BigQuery error:', error);
    res.status(500).send('An error occurred while recording the button tap');
  }
});

app.post('/analytics/:id/page-view', async (req, res) => {
  const { id } = req.params;
  const { pageName, timestamp, userType } = req.body;
  const datasetId = 'unischedule_analytics';
  const tableId = 'page_views';
  const rows = [{ pageName, timestamp, userType, userId: id }];
  try {
    await bigqueryClient.dataset(datasetId).table(tableId).insert(rows);
    res.status(201).send('Page view recorded');
  }
  catch (error) {
    console.error('BigQuery error:', error);
    res.status(500).send('An error occurred while recording the page view');
  }
});

app.post('/analytics/:id/event', async (req, res) => {
  const { id } = req.params;
  const { eventName, timestamp, userType } = req.body;
  const datasetId = 'unischedule_analytics';
  const tableId = 'events';
  const rows = [{ eventName, timestamp, userType, userId: id }];
  try {
    await bigqueryClient.dataset(datasetId).table(tableId).insert(rows);
    res.status(201).send('Event recorded');
  }
  catch (error) {
    console.error('BigQuery error:', error);
    res.status(500).send('An error occurred while recording the event');
  }
});

app.post('/analytics/:id/rating', async (req, res) => {
  const { id } = req.params;
  const { rating, classroom, timestamp, userType } = req.body;
  const datasetId = 'unischedule_analytics';
  const tableId = 'ratings';
  const rows = [{ rating, classroom, timestamp, userType, userId: id }];
  try {
    await bigqueryClient.dataset(datasetId).table(tableId).insert(rows);
    res.status(201).send('Rating recorded');
  }
  catch (error) {
    console.error('BigQuery error:', error);
    res.status(500).send('An error occurred while recording the rating');
  }
});

const port = process.env.PORT || 3000;

app.listen(port, () => {
console.log(`Server running at http://localhost:${port}`);
});
