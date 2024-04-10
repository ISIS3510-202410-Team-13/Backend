const express = require('express');
const admin = require('firebase-admin');
const { getFirestore } = require('firebase-admin/firestore');
const serviceAccount = require('../unischedule-5ee93-firebase-adminsdk-ci6y9-6cde344deb.json');

const db = getFirestore(admin.initializeApp({
    credential: admin.credential.cert(serviceAccount)
  }),'unischedule-backend');

const app = express();
app.use(express.json());

app.get('/user/:id/groups', async (req, res) => {
  const { id } = req.params;
  const userRef = await db.collection('Users').doc(id).get();

  if (!userRef.exists) {
    return res.status(404).send('User not found');
  }

  const userData = userRef.data();
  const groups = await Promise.all(userData.groups.map(async (groupId) => {
    const groupRef = await db.collection('Groups').doc(groupId).get();
    if (!groupRef.exists) return null;

    const groupData = groupRef.data();
    const memberProfiles = await Promise.all(groupData.members.slice(0, 3).map(async (memberId) => {
      const memberRef = await db.collection('Users').doc(memberId).get();
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
  const userRef = await db.collection('Users').doc(id).get();

  if (!userRef.exists) {
    return res.status(404).send('User not found');
  }

  const userData = userRef.data();
  const friends = await Promise.all(userData.friends.map(async (friendId) => {
    const friendRef = await db.collection('Users').doc(friendId).get();
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
  const userRef = await db.collection('Users').doc(id).get();

  if (!userRef.exists) {
    return res.status(404).send('User not found');
  }

  const userData = userRef.data();
  const events = await Promise.all(userData.events.map(async (eventId) => {
    const eventRef = await db.collection('Events').doc(eventId).get();
    return eventRef.exists ? { id: eventId, ...eventRef.data() } : null;
  }));

  res.status(200).json(events.filter(event => event !== null));
});

app.get('/group/:id/members', async (req, res) => {
    const { id } = req.params;
    const groupRef = await db.collection('Groups').doc(id).get();
  
    if (!groupRef.exists) {
      return res.status(404).send('Group not found');
    }
  
    const groupData = groupRef.data();
    const members = await Promise.all(groupData.members.map(async (memberId) => {
      const memberRef = await db.collection('Users').doc(memberId).get();
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
  const groupRef = await db.collection('Groups').doc(id).get();

  if (!groupRef.exists) {
    return res.status(404).send('Group not found');
  }

  const groupData = groupRef.data();
  const events = await Promise.all(groupData.events.map(async (eventId) => {
    const eventRef = await db.collection('Events').doc(eventId).get();
    return eventRef.exists ? { id: eventId, ...eventRef.data() } : null;
  }));

  res.status(200).json(events.filter(event => event !== null));
});

const port = process.env.PORT;

app.listen(port, () => {
console.log(`Server running at http://localhost:${port}`);
});

