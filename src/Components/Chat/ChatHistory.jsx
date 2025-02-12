import React from 'react';
import { List, ListItemButton, ListItemText, Divider, Box, Paper, Typography } from '@mui/material';

const ChatHistory = (currentUser, chats) => {
  return (
    <Box sx={{ p: 2, overflowY: 'auto', maxHeight: '400px' }}>
      {chats.map((m) => {
        const isSender = m.from === currentUser;
        const ts = new Date(m.timestamp * 1000).toLocaleString();

        return (
          <Paper
            key={m.id}
            sx={{
              p: 2,
              mb: 2,
              maxWidth: '75%',
              alignSelf: isSender ? 'flex-end' : 'flex-start',
              backgroundColor: isSender ? 'primary.main' : 'grey.300',
              color: isSender ? 'white' : 'black',
            }}
          >
            <Typography variant="body1">{m.message}</Typography>
            <Typography variant="caption" sx={{ display: 'block', textAlign: 'right', mt: 1 }}>
              {ts}
            </Typography>
          </Paper>
        );
      })}
    </Box>
  );
};

export default ChatHistory;