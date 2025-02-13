import React, { useState } from 'react';
import { useNavigate } from 'react-router-dom';
import {
  Container,
  Box,
  TextField,
  Button,
  Typography,
  Paper,
  Alert,
  InputAdornment,
} from '@mui/material';
import EditIcon from '@mui/icons-material/Edit';
import PersonIcon from '@mui/icons-material/Person';
import LockIcon from '@mui/icons-material/Lock';
import VpnKeyIcon from '@mui/icons-material/VpnKey';
import { ToastContainer, toast } from 'react-toastify';
import 'react-toastify/dist/ReactToastify.css';
import axiosInstance from '../utils/axiosInstance';

function Register() {
  const [credentials, setCredentials] = useState({
    username: '',
    password: '',
    confirmPassword: '',
  });
  const [error, setError] = useState('');
  const [passwordValid, setPasswordValid] = useState(true);
  const [passwordMatch, setPasswordMatch] = useState(true);
  const navigate = useNavigate();

  // Handle form field changes
  const handleChange = e => {
    setCredentials({ ...credentials, [e.target.name]: e.target.value });
  };

  // Password validation function
  const validatePassword = password => {
    // At least 8 characters, at least one letter, one number, and one special character
    const regex =
      /^(?=.*[A-Za-z])(?=.*\d)(?=.*[!@#$%^&*])[A-Za-z\d!@#$%^&*]{8,}$/;
    return regex.test(password);
  };

  // Handle form submission
  const handleSubmit = async e => {
    e.preventDefault();
    setError('');
    setPasswordValid(true);
    setPasswordMatch(true);

    // Check if password matches confirmPassword
    if (credentials.password !== credentials.confirmPassword) {
      setPasswordMatch(false);
      toast.error('Passwords do not match!');
      return;
    }

    // Validate password
    if (!validatePassword(credentials.password)) {
      setPasswordValid(false);
      toast.error(
        'Password must be at least 8 characters, contain a letter, a number, and a special character!'
      );
      return;
    }

    try {
      // const res = await axios.post(endpoint, credentials);
      const res = await axiosInstance.post('/register', credentials);
      if (res.data.status) {
        toast.success('Registration successful! Redirecting to chat...');
        navigate(`/chat?u=${credentials.username}`);
      } else {
        setError(res.data.message);
        toast.error(res.data.message);
      }
    } catch (error) {
      setError('Something went wrong. Please try again.');
      toast.error('Something went wrong. Please try again.');
    }
  };

  return (
    <Box
      sx={{
        backgroundColor: '#E7ECEF',
        display: 'flex',
        flexDirection: 'column',
        justifyContent: 'center',
        alignItems: 'center',
        minHeight: '90vh',
        p: 2,
      }}
    >
      <Container maxWidth="xs">
        <Paper
          elevation={12}
          sx={{
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
          }}
        >
          {/* Catchy Line */}
          <Typography
            variant="h5"
            fontWeight="bold"
            align="center"
            sx={{ mb: 3 }}
          >
            Register for ChatConnect
          </Typography>
          <Typography
            variant="body2"
            align="center"
            sx={{ color: '#A3CEF1', mb: 3 }}
          >
            Let's get you connected! Register now to start chatting.
          </Typography>

          {/* Error Alert */}
          {error && (
            <Alert
              severity="error"
              sx={{ mb: 2, backgroundColor: '#ff4d4d', color: '#ffffff' }}
            >
              {error}
            </Alert>
          )}

          {/* Username Field */}
          <TextField
            fullWidth
            label="Username"
            variant="outlined"
            name="username"
            value={credentials.username}
            onChange={handleChange}
            sx={{ mb: 2, backgroundColor: '#ffffff', borderRadius: 1 }}
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
            InputProps={{
              startAdornment: (
                <InputAdornment position="start">
                  <LockIcon sx={{ color: '#274C77' }} />
                </InputAdornment>
              ),
            }}
            helperText={
              passwordValid
                ? ''
                : 'Password must be at least 8 characters, contain a letter, a number, and a special character.'
            }
            error={!passwordValid}
          />

          {/* Confirm Password Field */}
          <TextField
            fullWidth
            label="Confirm Password"
            variant="outlined"
            type="password"
            name="confirmPassword"
            value={credentials.confirmPassword}
            onChange={handleChange}
            sx={{ mb: 2, backgroundColor: '#ffffff', borderRadius: 1 }}
            InputProps={{
              startAdornment: (
                <InputAdornment position="start">
                  <VpnKeyIcon sx={{ color: '#274C77' }} />
                </InputAdornment>
              ),
            }}
            helperText={passwordMatch ? '' : 'Passwords do not match.'}
            error={!passwordMatch}
          />

          {/* Register Button */}
          <Button
            fullWidth
            variant="contained"
            startIcon={<EditIcon />}
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
            Register
          </Button>
          {/* Existing user Section */}
          <Box sx={{ mt: 3, textAlign: 'center' }}>
            <Typography
              variant="body2"
              sx={{
                fontWeight: 'bold',
                color: '#A3CEF1',
                '&:hover': {
                  color: '#6096BA',
                  cursor: 'pointer',
                },
              }}
            >
              Have an Account?{' '}
              <span
                onClick={() => navigate('/login')}
                style={{ textDecoration: 'underline', color: '#e5383b' }}
              >
                Login
              </span>
            </Typography>
          </Box>
        </Paper>
      </Container>

      {/* Toast container to hold the messages */}
      <ToastContainer />
    </Box>
  );
}

export default Register;
