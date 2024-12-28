import { Typography } from '@mui/material';

/**
 * function is used to build the necessary object from the const steps to render the tour correctly.
 * @param {Array} staticSteps - an array of steps
 * @param {String} prefix - the prefix string to attach to the selector
 * @returns Array of steps with combined values of id, selector and content to build the tour properly.
 */
const derieveTourSteps = (staticSteps, prefix) => {
  return staticSteps.map(({ title }, index) => ({
    id: index,
    selector: `[data-tour="${prefix}-${index}"]`,
    content: <Typography variant="caption">{title}</Typography>,
  }));
};

const OVERVIEW_PAGE_TOUR_STEPS = [
  {
    title:
      'Welcome to Mashed Application Tutorial. Learn how to navigate the application in thorough detail. Each page is designed with its own help page, and you can go over it again and again as much as you would like.',
  },
  {
    title:
      'This is your asset summary. You can see how many assets are assigned to categories or plans. You can also see how many assets are due for return. Hover over the categories and plans to view the list containing said assets.',
  },
  {
    title:
      'This is your cost summary. You can view estimated cost of all the assets in your portfolio. Unestimated assets are the number of assets that do not have cost price associated with them.',
  },
  {
    title: 'Learn more about your asset categorization. Click on each content to navigate to the respective section.',
  },
  {
    title:
      'View the graphical representation of the combination of assets, categories and maintenance plans. Hover over them to see the count of each item.',
  },
  {
    title:
      'This is where you can see notifications for the list of maintenance plans that are due. Click on it to disable the notification.',
  },
  {
    title: 'Switch appearances or logout of revisit tutorial within the application from here.',
  },
];

const ASSETS_PAGE_TOUR_STEPS = [
  {
    title: 'This is where you can view the list of all your assets.',
  },
  {
    title: 'Search for any asset with ease from this button.',
  },
  { title: 'Change the view mode of the assets if you are not satisfied with the list layout' },
  {
    title:
      'Add assets at once, or one at a time from here. If you need to remove assets, select the ones you need to remove from the checkbox against each asset.',
  },
];

const CATEGORIES_PAGE_TOUR_STEPS = [
  {
    title: 'This is where you can view your categories.',
  },
  {
    title: 'Create new category from here',
  },
  {
    title: 'Download a report of the categories from here.',
  },
  {
    title: 'Filter your categories in terms of the status of the category. Eg, draft, archived etc',
  },
  {
    title: 'Sort categories from here. ',
  },
];

const MAINTENANCE_PAGE_TOUR_STEPS = [
  {
    title: 'This is where you can view your maintenance plans.',
  },
  {
    title: 'Create new maintenance plan from here',
  },
  {
    title: 'Download a report of the maintenance plans from here.',
  },
  {
    title: 'Filter your maintenance plans in terms of the status of the plan. Eg, draft, archived etc.',
  },
  {
    title: 'Sort maintenance plans from here. ',
  },
];

const REPORTS_PAGE_TOUR_STEPS = [
  {
    title:
      'View detailed reports of your assets from here. Reports are combined to instrument as a guide for you to understand your asset portfolio better. The data defaults to YTD to give you a perspective of how your assets perform over the course of the year.',
  },
  {
    title: 'Filter results to match a more reasonable timeline.',
  },
  {
    title: 'You can download or export the selected data from this button here.',
  },
  {
    title:
      'Displays the total valuation of all assets from the time selected. All assets are in terms of dollar value.',
  },
  {
    title:
      'Displays the valuation of assets that are categoried under either one of the categories. If the asset does not belong to any category, then it is exempt from this report section.',
  },
  {
    title: 'Displays the most recently added asset in this section.',
  },
  {
    title: 'Displays the most recent maintenance plan that was created.',
  },
  {
    title: 'View the list of assets that fall under the above filter results column.',
  },
];

const PROFILE_PAGE_TOUR_STEPS = [
  {
    title:
      'This is profile details where you can update details about yourself, choose your own avatar or view or edit settings.',
  },
  {
    title: 'Change avatar from here if you get bored with your profile picture.',
  },
  {
    title: 'This is your username and your biographic information.',
  },
  {
    title: 'Add a description about yourself if you would want',
  },
  {
    title: 'Displays your general preferences in the app',
  },
  {
    title: 'View general profile statistics such as the number of categories, plans or assets you have created.',
  },
  {
    title: 'Alter the appearance of the app from here.',
  },
];

const RECENT_ACTIVITIES_PAGE_TOUR_STEPS = [
  {
    title:
      'View a list of last 10 changes made with your assets here. If assets are shared between collaborators and change is made by another user, you can still view the changes here.',
  },
  {
    title: 'Download the history of all the changes made within the last 30 calendar days from today.',
  },
];

const PERSONAL_NOTES_PAGE_DATA_TOUR_STEPS = [
  {
    title: 'View list of notes that you have created over time.',
  },
  {
    title: 'Add new note from here. You can also assign location to the notes.',
  },
];

const SELECTED_CATEGORY_PAGE_DATA_TOUR_STEPS = [
  {
    title: 'View details about your selected category.',
  },
  {
    title: 'Mark category as favourites so that they can be displayed on the left menu bar easily.',
  },
  {
    title: 'Share the category and its items to other collaborators.',
  },
  {
    title: 'Upload images to represent the category',
  },
  {
    title:
      'Click on the add button to add assets to a selected category. Assets can be shared across many other categories.',
  },
  {
    title: 'Select assets that you would like to remove and press remove button.',
  },
  {
    title: 'View list of all assets that are assigned to this specific category.',
  },
  {
    title:
      'View a graph representation of the assets that belong to the category. If you have more than 10 items in the selected category, then only the last 10 are displayed.',
  },
];

const SELECTED_MAINTENANCE_PLAN_PAGE_DATA_TOUR_STEPS = [
  {
    title: 'View details about your selected maintenance plan.',
  },
  {
    title: 'Mark maintenance plan as favourites so that they can be displayed on the left menu bar easily.',
  },
  {
    title: 'Share the maintenance plan and its items to other collaborators.',
  },
  {
    title: 'Upload images to represent the maintenance plan.',
  },
  {
    title:
      'Click on the add button to add assets to a selected maintenance plan. Assets can be shared across many other maintenance plans.',
  },
  {
    title: 'Select assets that you would like to remove and press remove button.',
  },
  {
    title: 'View list of all assets that are assigned to this specific maintenance plan.',
  },
  {
    title:
      'View a graph representation of the assets that belong to the maintenance plan. If you have more than 10 items in the selected maintenance plan, then only the last 10 are displayed.',
  },
];

/**
 * Step Mapper function
 *
 * An object used to hold the steps and its respective start and end points. Based on the location of the
 * application, the step mapper will build a start and end number which will align to the tour steps.
 * The start number will be the length of previous tour and the end will always be the sum of the
 * previous and current tour since the end cannot be less than the start
 */
export const DEFAULT_STEP_MAPPER = {
  '/': {
    start: 0,
    end: OVERVIEW_PAGE_TOUR_STEPS.length,
  },
  '/inventories/list': {
    start: OVERVIEW_PAGE_TOUR_STEPS.length,
    end: OVERVIEW_PAGE_TOUR_STEPS.length + ASSETS_PAGE_TOUR_STEPS.length,
  },
  '/categories/list': {
    start: OVERVIEW_PAGE_TOUR_STEPS.length + ASSETS_PAGE_TOUR_STEPS.length,
    end: OVERVIEW_PAGE_TOUR_STEPS.length + ASSETS_PAGE_TOUR_STEPS.length + CATEGORIES_PAGE_TOUR_STEPS.length,
  },
  '/plans/list': {
    start: OVERVIEW_PAGE_TOUR_STEPS.length + ASSETS_PAGE_TOUR_STEPS.length + CATEGORIES_PAGE_TOUR_STEPS.length,
    end:
      OVERVIEW_PAGE_TOUR_STEPS.length +
      ASSETS_PAGE_TOUR_STEPS.length +
      CATEGORIES_PAGE_TOUR_STEPS.length +
      MAINTENANCE_PAGE_TOUR_STEPS.length,
  },
  '/reports': {
    start:
      OVERVIEW_PAGE_TOUR_STEPS.length +
      ASSETS_PAGE_TOUR_STEPS.length +
      CATEGORIES_PAGE_TOUR_STEPS.length +
      MAINTENANCE_PAGE_TOUR_STEPS.length,
    end:
      OVERVIEW_PAGE_TOUR_STEPS.length +
      ASSETS_PAGE_TOUR_STEPS.length +
      CATEGORIES_PAGE_TOUR_STEPS.length +
      MAINTENANCE_PAGE_TOUR_STEPS.length +
      REPORTS_PAGE_TOUR_STEPS.length,
  },
  '/profile': {
    start:
      OVERVIEW_PAGE_TOUR_STEPS.length +
      ASSETS_PAGE_TOUR_STEPS.length +
      CATEGORIES_PAGE_TOUR_STEPS.length +
      MAINTENANCE_PAGE_TOUR_STEPS.length +
      REPORTS_PAGE_TOUR_STEPS.length,
    end:
      OVERVIEW_PAGE_TOUR_STEPS.length +
      ASSETS_PAGE_TOUR_STEPS.length +
      CATEGORIES_PAGE_TOUR_STEPS.length +
      MAINTENANCE_PAGE_TOUR_STEPS.length +
      REPORTS_PAGE_TOUR_STEPS.length +
      PROFILE_PAGE_TOUR_STEPS.length,
  },
  '/recent/activities': {
    start:
      OVERVIEW_PAGE_TOUR_STEPS.length +
      ASSETS_PAGE_TOUR_STEPS.length +
      CATEGORIES_PAGE_TOUR_STEPS.length +
      MAINTENANCE_PAGE_TOUR_STEPS.length +
      REPORTS_PAGE_TOUR_STEPS.length +
      PROFILE_PAGE_TOUR_STEPS.length,
    end:
      OVERVIEW_PAGE_TOUR_STEPS.length +
      ASSETS_PAGE_TOUR_STEPS.length +
      CATEGORIES_PAGE_TOUR_STEPS.length +
      MAINTENANCE_PAGE_TOUR_STEPS.length +
      REPORTS_PAGE_TOUR_STEPS.length +
      PROFILE_PAGE_TOUR_STEPS.length +
      RECENT_ACTIVITIES_PAGE_TOUR_STEPS.length,
  },
  '/profile/notes': {
    start:
      OVERVIEW_PAGE_TOUR_STEPS.length +
      ASSETS_PAGE_TOUR_STEPS.length +
      CATEGORIES_PAGE_TOUR_STEPS.length +
      MAINTENANCE_PAGE_TOUR_STEPS.length +
      REPORTS_PAGE_TOUR_STEPS.length +
      PROFILE_PAGE_TOUR_STEPS.length +
      RECENT_ACTIVITIES_PAGE_TOUR_STEPS.length,
    end:
      OVERVIEW_PAGE_TOUR_STEPS.length +
      ASSETS_PAGE_TOUR_STEPS.length +
      CATEGORIES_PAGE_TOUR_STEPS.length +
      MAINTENANCE_PAGE_TOUR_STEPS.length +
      REPORTS_PAGE_TOUR_STEPS.length +
      PROFILE_PAGE_TOUR_STEPS.length +
      RECENT_ACTIVITIES_PAGE_TOUR_STEPS.length +
      PERSONAL_NOTES_PAGE_DATA_TOUR_STEPS.length,
  },
  '/category/id': {
    start:
      OVERVIEW_PAGE_TOUR_STEPS.length +
      ASSETS_PAGE_TOUR_STEPS.length +
      CATEGORIES_PAGE_TOUR_STEPS.length +
      MAINTENANCE_PAGE_TOUR_STEPS.length +
      REPORTS_PAGE_TOUR_STEPS.length +
      PROFILE_PAGE_TOUR_STEPS.length +
      RECENT_ACTIVITIES_PAGE_TOUR_STEPS.length +
      PERSONAL_NOTES_PAGE_DATA_TOUR_STEPS.length,
    end:
      OVERVIEW_PAGE_TOUR_STEPS.length +
      ASSETS_PAGE_TOUR_STEPS.length +
      CATEGORIES_PAGE_TOUR_STEPS.length +
      MAINTENANCE_PAGE_TOUR_STEPS.length +
      REPORTS_PAGE_TOUR_STEPS.length +
      PROFILE_PAGE_TOUR_STEPS.length +
      RECENT_ACTIVITIES_PAGE_TOUR_STEPS.length +
      PERSONAL_NOTES_PAGE_DATA_TOUR_STEPS.length +
      SELECTED_CATEGORY_PAGE_DATA_TOUR_STEPS.length,
  },
  '/plan/id': {
    start:
      OVERVIEW_PAGE_TOUR_STEPS.length +
      ASSETS_PAGE_TOUR_STEPS.length +
      CATEGORIES_PAGE_TOUR_STEPS.length +
      MAINTENANCE_PAGE_TOUR_STEPS.length +
      REPORTS_PAGE_TOUR_STEPS.length +
      PROFILE_PAGE_TOUR_STEPS.length +
      RECENT_ACTIVITIES_PAGE_TOUR_STEPS.length +
      PERSONAL_NOTES_PAGE_DATA_TOUR_STEPS.length +
      SELECTED_CATEGORY_PAGE_DATA_TOUR_STEPS.length,
    end:
      OVERVIEW_PAGE_TOUR_STEPS.length +
      ASSETS_PAGE_TOUR_STEPS.length +
      CATEGORIES_PAGE_TOUR_STEPS.length +
      MAINTENANCE_PAGE_TOUR_STEPS.length +
      REPORTS_PAGE_TOUR_STEPS.length +
      PROFILE_PAGE_TOUR_STEPS.length +
      RECENT_ACTIVITIES_PAGE_TOUR_STEPS.length +
      PERSONAL_NOTES_PAGE_DATA_TOUR_STEPS.length +
      SELECTED_CATEGORY_PAGE_DATA_TOUR_STEPS.length +
      SELECTED_MAINTENANCE_PLAN_PAGE_DATA_TOUR_STEPS.length,
  },
};

const DEFAULT_TOUR_STEPS = [
  ...derieveTourSteps(OVERVIEW_PAGE_TOUR_STEPS, 'overview'),
  ...derieveTourSteps(ASSETS_PAGE_TOUR_STEPS, 'assets'),
  ...derieveTourSteps(CATEGORIES_PAGE_TOUR_STEPS, 'categories'),
  ...derieveTourSteps(MAINTENANCE_PAGE_TOUR_STEPS, 'plans'),
  ...derieveTourSteps(REPORTS_PAGE_TOUR_STEPS, 'reports'),
  ...derieveTourSteps(PROFILE_PAGE_TOUR_STEPS, 'profile'),
  ...derieveTourSteps(RECENT_ACTIVITIES_PAGE_TOUR_STEPS, 'recentActivities'),
  ...derieveTourSteps(PERSONAL_NOTES_PAGE_DATA_TOUR_STEPS, 'notes'),
  ...derieveTourSteps(SELECTED_CATEGORY_PAGE_DATA_TOUR_STEPS, 'selected-category'),
  ...derieveTourSteps(SELECTED_MAINTENANCE_PLAN_PAGE_DATA_TOUR_STEPS, 'selected-plan'),
];

export default DEFAULT_TOUR_STEPS;
