import React from 'react';
import { Navigate } from 'react-router-dom';
import { useAuth } from '../contexts/AuthContext';
import { Box, CircularProgress, Typography, Alert, Button } from '@mui/material';

interface AdminRouteProps {
  children: React.ReactNode;
}

const AdminRoute: React.FC<AdminRouteProps> = ({ children }) => {
  const { user, loading } = useAuth();

  if (loading) {
    return (
      <Box
        display="flex"
        flexDirection="column"
        justifyContent="center"
        alignItems="center"
        minHeight="100vh"
      >
        <CircularProgress />
        <Typography variant="body2" sx={{ mt: 2 }}>
          验证管理员权限中...
        </Typography>
      </Box>
    );
  }

  if (!user) {
    return <Navigate to="/login" replace />;
  }

  // 检查用户是否为管理员
  if (user.role !== 'admin') {
    return (
      <Box
        display="flex"
        flexDirection="column"
        justifyContent="center"
        alignItems="center"
        minHeight="100vh"
        sx={{ p: 3 }}
      >
        <Alert severity="error" sx={{ mb: 2, maxWidth: 400 }}>
          <Typography variant="h6" gutterBottom>
            访问被拒绝
          </Typography>
          <Typography variant="body2" sx={{ mb: 2 }}>
            您没有管理员权限，无法访问此页面。
          </Typography>
          <Button
            variant="contained"
            onClick={() => window.history.back()}
            sx={{ mt: 1 }}
          >
            返回上一页
          </Button>
        </Alert>
      </Box>
    );
  }

  return <>{children}</>;
};

export default AdminRoute;
