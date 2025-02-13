import React, { Component } from 'react';
import SocketConnection from '../../socket-connection';
import ChatHistory from './ChatHistory';
import ContactList from './ContactList';
import axiosInstance from '../../utils/axiosInstance';
import {
  Container,
  TextField,
  Box,
  Button,
  Typography,
  Paper,
  List,
  ListItem,
  ListItemText,
} from '@mui/material';
import SendIcon from '@mui/icons-material/Send';
import AddIcon from '@mui/icons-material/Add';

class Chat extends Component {
  constructor(props) {
    super(props);
    this.state = {
      socketConn: '',
      username: '',
      message: '',
      to: '',
      isInvalid: false,
      contact: '',
      contacts: [],
      renderContactList: [],
      chats: [],
      chatHistory: [],
    };
  }

  componentDidMount = async () => {
    const queryParams = new URLSearchParams(window.location.search);
    const user = queryParams.get('u');
    this.setState({ username: user });
    this.getContacts(user);

    const conn = new SocketConnection();
    await this.setState({ socketConn: conn });
    this.state.socketConn.connect(message => {
      const msg = JSON.parse(message.data);
      if (this.state.to === msg.from || this.state.username === msg.from) {
        this.setState(
          { chats: [...this.state.chats, msg] },
          () => this.renderChatHistory(this.state.username, this.state.chats)
        );
      }
    });
    this.state.socketConn.connected(user);
  };

  onChange = event => {
    this.setState({ [event.target.name]: event.target.value });
  };

  onSubmit = e => {
    e.preventDefault();
  
    // if (!this.state.to || !this.state.message) {
    //   console.log('Recipient or message is empty');
    //   return;
    // }
  
    // Create the message object with the required format
    const msg = {
      type: 'chat',
      chat: {
        from: this.state.username,  // from the logged-in user
        to: this.state.to,          // to the selected contact
        message: this.state.message, // the message content
      },
    };
  
    if (this.state.socketConn && this.state.socketConn.socket.readyState === WebSocket.OPEN) {
      this.state.socketConn.sendMsg(msg);
      this.setState({ message: '' });
    } else {
      console.log('WebSocket not connected. Message not sent.');
    }
  };
  
  

  getContacts = async user => {
    const res = await axiosInstance.get(`/contact-list?username=${user}`);
    if (res.data['data']) {
      this.setState({ contacts: res.data.data });
      this.renderContactList(res.data.data);
    }
  };

  fetchChatHistory = async (u1, u2) => {
    try {
      const res = await axiosInstance.get(`/chat-history?u1=${u1}&u2=${u2}`);;
      if (res.data.status && Array.isArray(res.data.data)) {
        this.setState({ chats: res.data.data.reverse() });
        this.renderChatHistory(u1, res.data.data);
      } else {
        this.setState({ chatHistory: [] });
      }
    } catch (error) {
      console.error('Error fetching chat history:', error);
      this.setState({ chatHistory: [] });
    }
  };

  renderChatHistory = (currentUser, chats) => {
    this.setState({ chatHistory: ChatHistory(currentUser, chats) });
  };

  renderContactList = contacts => {
    this.setState({ renderContactList: ContactList(contacts, this.sendMessageTo) });
  };

  sendMessageTo = to => {
    // Set the recipient ("to") user and fetch the chat history
    this.setState({ to });
    this.fetchChatHistory(this.state.username, to);
    this.state.socketConn.sendAckMessage(this.state.username)
  };
  

  render() {
    return (
      <Container maxWidth="md" sx={{ mt: 4, p: 2, display: 'flex', flexDirection: 'column' }}>
        <Typography variant="h6" align="right" gutterBottom>{this.state.username}</Typography>
        
        <Paper sx={{ p: 2, mb: 2, display: 'flex', gap: 1 }}>
          <TextField 
            fullWidth 
            label="Add Contact" 
            variant="outlined" 
            name="contact" 
            value={this.state.contact} 
            onChange={this.onChange} 
          />
          <Button variant="contained" color="primary" startIcon={<AddIcon />} onClick={this.addContact}>
            Add
          </Button>
        </Paper>
        
        <Box sx={{ display: 'flex', height: '65vh', gap: 2 }}>
          <Paper sx={{ width: '30%', overflowY: 'auto', p: 1 }}>
            <List>
              {this.state.contacts.map(contact => (
                <ListItem button key={contact.username} onClick={() => this.sendMessageTo(contact.username)}>
                  <ListItemText primary={contact.username} />
                </ListItem>
              ))}
            </List>
          </Paper>
          
          <Paper sx={{ flexGrow: 1, display: 'flex', flexDirection: 'column' }}>
            <Box sx={{ flexGrow: 1, overflowY: 'auto', p: 2, display: 'flex', flexDirection: 'column-reverse' }}>
              {this.state.chatHistory}
            </Box>
            <Box sx={{ p: 2, borderTop: '1px solid #ddd' }}>
              <TextField 
                fullWidth 
                multiline 
                minRows={2} 
                placeholder="Type your message..." 
                name="message" 
                value={this.state.message} 
                onChange={this.onChange} 
                // onKeyDown={this.onSubmit} 
              />
              <Button fullWidth variant="contained" color="primary" startIcon={<SendIcon />} onClick={this.onSubmit}>
                Send
              </Button>
            </Box>
          </Paper>
        </Box>
      </Container>
    );
  }
}

export default Chat;
