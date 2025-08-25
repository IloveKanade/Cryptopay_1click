import React from 'react';
import {
  AppBar,
  Toolbar,
  Typography,
  Button,
  Box,
  IconButton,
  Drawer,
  List,
  ListItem,
  ListItemIcon,
  ListItemText,
  Divider,
  useTheme,
  useMediaQuery,
} from '@mui/material';
import {
  Menu as MenuIcon,
  Dashboard as DashboardIcon,
  People as PeopleIcon,
  Receipt as OrdersIcon,
  AccountBalanceWallet as WalletIcon,
  Settings as SettingsIcon,
  Logout as LogoutIcon,
} from '@mui/icons-material';
import { useNavigate, useLocation } from 'react-router-dom';
import { useAuth } from '../contexts/AuthContext';

const AdminNavigation: React.FC = () => {
  const [drawerOpen, setDrawerOpen] = React.useState(false);
  const { logout } = useAuth();
  const navigate = useNavigate();
  const location = useLocation();
  const theme = useTheme();
  const isMobile = useMediaQuery(theme.breakpoints.down('md'));

  const menuItems = [
    { text: '仪表板', icon: <DashboardIcon />, path: '/admin' },
    { text: '用户管理', icon: <PeopleIcon />, path: '/admin/users' },
    { text: '订单管理', icon: <OrdersIcon />, path: '/admin/orders' },
    { text: '中间人钱包', icon: <WalletIcon />, path: '/admin/middleman' },
    { text: '系统设置', icon: <SettingsIcon />, path: '/admin/settings' },
  ];

  const handleLogout = () => {
    logout();
    navigate('/login');
  };

  const handleNavigate = (path: string) => {
    navigate(path);
    if (isMobile) {
      setDrawerOpen(false);
    }
  };

  const drawer = (
    <Box sx={{ width: 250 }}>
      <Box sx={{ p: 2, borderBottom: '1px solid rgba(255, 255, 255, 0.1)' }}>
        <Typography variant="h6" sx={{ color: '#00ff88', fontWeight: 'bold' }}>
          管理员控制台
        </Typography>
      </Box>
      <List>
        {menuItems.map((item) => (
          <ListItem
            key={item.path}
            button
            onClick={() => handleNavigate(item.path)}
            sx={{
              backgroundColor: location.pathname === item.path ? 'rgba(0, 255, 136, 0.1)' : 'transparent',
              borderLeft: location.pathname === item.path ? '4px solid #00ff88' : 'none',
              '&:hover': {
                backgroundColor: 'rgba(0, 255, 136, 0.05)',
              },
            }}
          >
            <ListItemIcon sx={{ color: location.pathname === item.path ? '#00ff88' : '#b0b0b0' }}>
              {item.icon}
            </ListItemIcon>
            <ListItemText
              primary={item.text}
              sx={{
                color: location.pathname === item.path ? '#00ff88' : '#ffffff',
                fontWeight: location.pathname === item.path ? 'bold' : 'normal',
              }}
            />
          </ListItem>
        ))}
      </List>
      <Divider sx={{ borderColor: 'rgba(255, 255, 255, 0.1)' }} />
      <List>
        <ListItem button onClick={handleLogout}>
          <ListItemIcon sx={{ color: '#ff4444' }}>
            <LogoutIcon />
          </ListItemIcon>
          <ListItemText primary="退出登录" sx={{ color: '#ff4444' }} />
        </ListItem>
      </List>
    </Box>
  );

  return (
    <>
      <AppBar
        position="static"
        sx={{
          background: 'rgba(26, 26, 46, 0.95)',
          backdropFilter: 'blur(10px)',
          borderBottom: '1px solid rgba(0, 255, 136, 0.2)',
        }}
      >
        <Toolbar>
          {isMobile && (
            <IconButton
              edge="start"
              color="inherit"
              onClick={() => setDrawerOpen(true)}
              sx={{ mr: 2 }}
            >
              <MenuIcon />
            </IconButton>
          )}
          
          <Typography variant="h6" component="div" sx={{ flexGrow: 1, color: '#00ff88', fontWeight: 'bold' }}>
            管理员控制台
          </Typography>

          {!isMobile && (
            <Box sx={{ display: 'flex', gap: 1 }}>
              {menuItems.map((item) => (
                <Button
                  key={item.path}
                  color="inherit"
                  onClick={() => handleNavigate(item.path)}
                  sx={{
                    color: location.pathname === item.path ? '#00ff88' : '#ffffff',
                    fontWeight: location.pathname === item.path ? 'bold' : 'normal',
                    borderBottom: location.pathname === item.path ? '2px solid #00ff88' : 'none',
                    '&:hover': {
                      backgroundColor: 'rgba(0, 255, 136, 0.1)',
                    },
                  }}
                >
                  {item.text}
                </Button>
              ))}
              <Button
                color="inherit"
                onClick={handleLogout}
                sx={{
                  color: '#ff4444',
                  '&:hover': {
                    backgroundColor: 'rgba(255, 68, 68, 0.1)',
                  },
                }}
              >
                退出登录
              </Button>
            </Box>
          )}
        </Toolbar>
      </AppBar>

      <Drawer
        anchor="left"
        open={drawerOpen}
        onClose={() => setDrawerOpen(false)}
        sx={{
          '& .MuiDrawer-paper': {
            background: 'rgba(26, 26, 46, 0.95)',
            backdropFilter: 'blur(10px)',
            borderRight: '1px solid rgba(0, 255, 136, 0.2)',
          },
        }}
      >
        {drawer}
      </Drawer>
    </>
  );
};

export default AdminNavigation;
