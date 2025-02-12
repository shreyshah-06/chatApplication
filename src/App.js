import React from 'react';
import { BrowserRouter, Routes, Route } from 'react-router-dom';
import { ThemeProvider, CssBaseline, Container } from '@mui/material';
import theme from './theme';
import Header from './Components/Header';
import Landing from './Components/Landing';
import Register from './Components/Register';
import Login from './Components/Login';
import Chat from './Components/Chat/Chat';
import Footer from './Components/Footer';

function App() {
  return (
    <ThemeProvider theme={theme}>
      <CssBaseline />
      <Container maxWidth="xl" disableGutters>
        <BrowserRouter>
          <Header />
          <Routes>
            <Route path="/" element={<Landing />} />
            <Route path="/register" element={<Register />} />
            <Route path="/login" element={<Login />} />
            <Route path="/chat" element={<Chat />} />
          </Routes>
          <Footer />
        </BrowserRouter>
      </Container>
    </ThemeProvider>
  );
}

export default App;
