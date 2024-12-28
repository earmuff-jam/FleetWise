import Popper from '@material-ui/core/Popper/Popper';
import { ArrowDropDownCircleRounded, DarkModeOutlined, LightModeOutlined, LogoutRounded } from '@mui/icons-material';
import {
  Button,
  ButtonGroup,
  ClickAwayListener,
  Grow,
  MenuItem,
  MenuList,
  Paper,
  Stack,
  Typography,
} from '@mui/material';

import { useRef, useState } from 'react';
import { useTour } from '@reactour/tour';
import { useLocation, useParams } from 'react-router-dom';
import { authActions } from '@features/LandingPage/authSlice';
import { useDispatch } from 'react-redux';
import { profileActions } from '@features/Profile/profileSlice';
import DEFAULT_TOUR_STEPS, { DEFAULT_STEP_MAPPER } from '@utils/tour/steps';

export default function AppToolBarMoreButtons({ profileDetails }) {
  const { id } = useParams();
  const location = useLocation();
  const dispatch = useDispatch();

  const { setIsOpen, setCurrentStep, setSteps } = useTour();

  const anchorRef = useRef(null);
  const [open, setOpen] = useState(false);

  const setTour = () => {
    let uri = location.pathname;
    if (id && uri.includes('/category/')) {
      uri = '/category/id';
    }
    if (id && uri.includes('/plan/')) {
      uri = '/plan/id';
    }

    const currentStep = DEFAULT_STEP_MAPPER[uri];
    const formattedSteps = DEFAULT_TOUR_STEPS.slice(currentStep.start, currentStep.end);
    setIsOpen(true);
    setCurrentStep(0);
    setSteps(formattedSteps);
  };

  const handleLogout = () => {
    dispatch(authActions.getLogout());
    localStorage.clear();
    window.location.href = '/';
  };

  const handleAppearance = () => {
    const draftData = { ...profileDetails, appearance: !profileDetails.appearance || false };
    dispatch(profileActions.updateProfileDetails({ draftData }));
  };

  const DEFAULT_MENU_OPTIONS = [
    {
      id: 1,
      label: 'Log off',
      icon: <LogoutRounded color="primary" />,
      action: handleLogout,
    },
    {
      id: 2,
      label: 'Change appearance',
      icon: profileDetails?.appearance ? (
        <LightModeOutlined color="text.secondary" fontSize="small" />
      ) : (
        <DarkModeOutlined fontSize="small" color="text.secondary" />
      ),
      action: handleAppearance,
    },
  ];

  const handleToggle = () => {
    setOpen((prevOpen) => !prevOpen);
  };

  const handleClose = (event) => {
    if (anchorRef.current && anchorRef.current.contains(event.target)) {
      return;
    }

    setOpen(false);
  };

  return (
    <>
      <ButtonGroup
        variant="contained"
        ref={anchorRef}
        aria-label="Button group with a nested menu"
        data-tour="overview-6"
      >
        <Button onClick={setTour}>Help</Button>
        <Button
          size="small"
          aria-controls={open ? 'split-button-menu' : undefined}
          aria-expanded={open ? 'true' : undefined}
          aria-label="select merge strategy"
          aria-haspopup="menu"
          onClick={handleToggle}
        >
          <ArrowDropDownCircleRounded />
        </Button>
      </ButtonGroup>
      <Popper sx={{ zIndex: 1 }} open={open} anchorEl={anchorRef.current} role={undefined} transition disablePortal>
        {({ TransitionProps, placement }) => (
          <Grow
            {...TransitionProps}
            style={{
              transformOrigin: placement === 'bottom' ? 'center top' : 'center bottom',
            }}
          >
            <Paper>
              <ClickAwayListener onClickAway={handleClose}>
                <MenuList id="split-button-menu" autoFocusItem>
                  {DEFAULT_MENU_OPTIONS.map((option) => (
                    <MenuItem
                      key={option.id}
                      onClick={() => {
                        option.action();
                        handleToggle();
                      }}
                    >
                      <Stack direction="row" spacing={1} alignItems="center">
                        {option.icon}
                        <Typography variant="caption">{option.label}</Typography>
                      </Stack>
                    </MenuItem>
                  ))}
                </MenuList>
              </ClickAwayListener>
            </Paper>
          </Grow>
        )}
      </Popper>
    </>
  );
}
