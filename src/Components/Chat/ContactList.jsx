import React from 'react';
import { List, ListItemButton, ListItemText, Divider, Box, Paper, Typography } from '@mui/material';

const ContactList = (contacts, sendMessage) => {
    return (
      <List>
        {contacts.map((c) => (
          <React.Fragment key={c.username}>
            <ListItemButton onClick={() => sendMessage(c.username)}>
              <ListItemText
                primary={c.username}
                secondary={`Last active: ${new Date(c.last_activity * 1000).toLocaleDateString()}`}
              />
            </ListItemButton>
            <Divider />
          </React.Fragment>
        ))}
      </List>
    );
  };
  

export default ContactList;
