import React from 'react';
import { Link } from 'react-router-dom';
import { Box, Button, Container, Typography, Grid, Paper } from '@mui/material';
import ChatIcon from '@mui/icons-material/Chat';
import EditIcon from '@mui/icons-material/Edit';
import ArrowForwardIcon from '@mui/icons-material/ArrowForward';

function Landing() {
  return (
    <Box sx={{ backgroundColor: '#E7ECEF', py: 8 }}>
      <Container maxWidth="xl">
        {/* Hero Section */}
        <Paper elevation={4} sx={{ p: 6, borderRadius: 3, backgroundColor: '#274C77', color: '#ffffff', alignItems: 'center' }}>
          <Box display="flex" justifyContent="center" alignItems="center" mb={4}>
            <ChatIcon sx={{ fontSize: 70, mr: 2, color: '#A3CEF1' }} />
            <Typography variant="h2" fontWeight="bold">
              ChatConnect
            </Typography>
          </Box>
          <Typography variant="h5" gutterBottom sx={{ color: '#A3CEF1', mb: 4, textAlign: 'center' }}>
            Connect seamlessly with people around the world in real-time.
          </Typography>

          <Grid container justifyContent="center" spacing={4}>
            <Grid item>
              <Link to="register" style={{ textDecoration: 'none' }}>
                <Button
                  size="large"
                  startIcon={<EditIcon />}
                  variant="contained"
                  sx={{
                    backgroundColor: '#6096BA',
                    '&:hover': {
                      backgroundColor: '#4a78a2', 
                      transform: 'scale(1.05)',
                    },
                  }}
                >
                  Register
                </Button>
              </Link>
            </Grid>
            <Grid item>
              <Link to="login" style={{ textDecoration: 'none' }}>
                <Button
                  size="large"
                  endIcon={<ArrowForwardIcon />}
                  variant="outlined"
                  sx={{
                    color: '#A3CEF1',
                    borderColor: '#A3CEF1',
                    '&:hover': {
                      borderColor: '#6096BA',
                      backgroundColor: '#A3CEF1', 
                      color: '#274C77', 
                    },
                  }}
                >
                  Login
                </Button>
              </Link>
            </Grid>
          </Grid>
        </Paper>

        {/* Features Section */}
        <Box sx={{ mt: 8 }}>
          <Typography variant="h4" align="center" color="#274C77" fontWeight="bold">
            Why Choose ChatConnect?
          </Typography>
          <Grid container spacing={4} sx={{ mt: 4 }}>
            <Grid item xs={12} md={4}>
              <Paper
                elevation={3}
                sx={{
                  p: 4,
                  backgroundColor: '#6096BA',
                  color: '#ffffff',
                  borderRadius: 3,
                  transition: 'all 0.3s ease-in-out', // Smooth transition
                  '&:hover': {
                    transform: 'scale(1.05)', // Scale effect on hover
                    boxShadow: '0 8px 16px rgba(0, 0, 0, 0.2)', // Shadow effect on hover
                  },
                }}
              >
                <Typography variant="h6">Real-Time Messaging</Typography>
                <Typography variant="body2" sx={{ mt: 1 }}>
                  Enjoy seamless, fast communication without any delays.
                </Typography>
              </Paper>
            </Grid>
            <Grid item xs={12} md={4}>
              <Paper
                elevation={3}
                sx={{
                  p: 4,
                  backgroundColor: '#8B8C89',
                  color: '#ffffff',
                  borderRadius: 3,
                  transition: 'all 0.3s ease-in-out', // Smooth transition
                  '&:hover': {
                    transform: 'scale(1.05)', // Scale effect on hover
                    boxShadow: '0 8px 16px rgba(0, 0, 0, 0.2)', // Shadow effect on hover
                  },
                }}
              >
                <Typography variant="h6">Secure Conversations</Typography>
                <Typography variant="body2" sx={{ mt: 1 }}>
                  All your chats are encrypted for maximum privacy.
                </Typography>
              </Paper>
            </Grid>
            <Grid item xs={12} md={4}>
              <Paper
                elevation={3}
                sx={{
                  p: 4,
                  backgroundColor: '#A3CEF1',
                  color: '#274C77',
                  borderRadius: 3,
                  transition: 'all 0.3s ease-in-out', // Smooth transition
                  '&:hover': {
                    transform: 'scale(1.05)', // Scale effect on hover
                    boxShadow: '0 8px 16px rgba(0, 0, 0, 0.2)', // Shadow effect on hover
                  },
                }}
              >
                <Typography variant="h6">Cross-Platform Support</Typography>
                <Typography variant="body2" sx={{ mt: 1 }}>
                  Access your chats on any device, anywhere.
                </Typography>
              </Paper>
            </Grid>
          </Grid>
        </Box>
      </Container>
    </Box>
  );
}

export default Landing;
