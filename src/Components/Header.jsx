import React from 'react';
import { AppBar, Toolbar, Typography, Container } from '@mui/material';
import ChatIcon from '@mui/icons-material/Chat';
import { Link } from 'react-router-dom';

function Header() {
  return (
    <AppBar position="static" sx={{ backgroundColor: '#6096BA' }}>
      <Container maxWidth="100%">
        <Toolbar>
          <Link to="/" style={{ textDecoration: 'none', color: 'inherit' }}>
            <Typography variant="h4" fontWeight="bold" sx={{ display: 'flex', alignItems: 'center' }}>
              <ChatIcon sx={{ fontSize: 40, mr: 1 }} />
              ChatConnect
            </Typography>
          </Link>
        </Toolbar>
      </Container>
    </AppBar>
  );
}

export default Header;
