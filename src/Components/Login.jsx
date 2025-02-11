import React, { useState } from 'react';
import { useNavigate } from 'react-router-dom';
import axios from 'axios';
import { Container, Box, TextField, Button, Typography, Paper, Alert, InputAdornment } from '@mui/material';
import LoginIcon from '@mui/icons-material/Login';
import PersonIcon from '@mui/icons-material/Person';
import LockIcon from '@mui/icons-material/Lock';
import ChatIcon from '@mui/icons-material/Chat';

function Login() {
  const [credentials, setCredentials] = useState({ username: '', password: '' });
  const [error, setError] = useState('');
  const navigate = useNavigate();
  const endpoint = 'http://localhost:8080/login';

  const handleChange = (e) => {
    setCredentials({ ...credentials, [e.target.name]: e.target.value });
  };

  const handleSubmit = async (e) => {
    e.preventDefault();
    setError('');
    try {
      const res = await axios.post(endpoint, credentials);
      if (res.data.status) {
        navigate(`/chat?u=${credentials.username}`);
      } else {
        setError(res.data.message);
      }
    } catch (error) {
      setError('Something went wrong. Please try again.');
    }
  };

  return (
    <Box sx={{ backgroundColor: '#E7ECEF', display: 'flex', justifyContent: 'center', alignItems: 'center', minHeight: '100vh', p: 2 }}>
      <Container maxWidth="xs">
        <Paper elevation={12} sx={{
          p: 4,
          borderRadius: 3,
          background: 'linear-gradient(135deg, #6096BA 0%, #274C77 100%)',
          color: '#ffffff',
          boxShadow: '0 4px 20px rgba(0, 0, 0, 0.1)',
          transform: 'translateY(-20px)',
          transition: 'transform 0.3s ease-in-out',
          '&:hover': {
            transform: 'translateY(-5px)',
          },
        }}>
          {/* Header */}
          <Box display="flex" justifyContent="center" alignItems="center" mb={2}>
            <ChatIcon sx={{ fontSize: 40, mr: 1, color: '#A3CEF1' }} />
            <Typography variant="h5" fontWeight="bold">
              ChatConnect
            </Typography>
          </Box>
          <Typography variant="body1" align="center" sx={{ color: '#A3CEF1', mb: 3 }}>
            Connect with the world in real-time. Login to get started.
          </Typography>

          {/* Error Alert */}
          {error && <Alert severity="error" sx={{ mb: 2, backgroundColor: '#ff4d4d', color: '#ffffff' }}>{error}</Alert>}

          {/* Username Field */}
          <TextField
            fullWidth
            label="Username"
            variant="outlined"
            name="username"
            value={credentials.username}
            onChange={handleChange}
            sx={{ mb: 2, backgroundColor: '#ffffff', borderRadius: 1 }}
            InputLabelProps={{ shrink: true }}
            InputProps={{
              startAdornment: (
                <InputAdornment position="start">
                  <PersonIcon sx={{ color: '#274C77' }} />
                </InputAdornment>
              ),
            }}
          />

          {/* Password Field */}
          <TextField
            fullWidth
            label="Password"
            variant="outlined"
            type="password"
            name="password"
            value={credentials.password}
            onChange={handleChange}
            sx={{ mb: 2, backgroundColor: '#ffffff', borderRadius: 1 }}
            InputLabelProps={{ shrink: true }}
            InputProps={{
              startAdornment: (
                <InputAdornment position="start">
                  <LockIcon sx={{ color: '#274C77' }} />
                </InputAdornment>
              ),
            }}
          />

          {/* Submit Button */}
          <Button
            fullWidth
            variant="contained"
            startIcon={<LoginIcon />}
            sx={{
              backgroundColor: '#6096BA',
              color: '#ffffff',
              padding: '10px 20px',
              fontSize: '16px',
              '&:hover': {
                backgroundColor: '#4a78a2',
                boxShadow: '0 8px 16px rgba(0, 0, 0, 0.2)',
              },
            }}
            onClick={handleSubmit}
          >
            Login
          </Button>
        </Paper>
      </Container>
    </Box>
  );
}

export default Login;
