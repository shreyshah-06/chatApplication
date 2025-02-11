import React from 'react';
import { Box, Typography, Container } from '@mui/material';

function Footer() {
  return (
    <Box
      component="footer"
      sx={{
        backgroundColor: '#284b63',
        color: '#ffffff',
        py: 2,
        position: 'relative',
        bottom: 0,
        width: '100%',
      }}
    >
      <Container>
        <Typography variant="body2" align="center">
          Powered By ChatConnect. All rights reserved.
        </Typography>
      </Container>
    </Box>
  );
}

export default Footer;
