export const paths = {
  home: {
    path: '/',
    getHref: () => '/',
  },

  users: {
    path: 'users/:gameName/:tagLine',
    getHref: (gameName: string, tagLine: string) =>
      `/users/${gameName}/${tagLine}`,
  },
} as const;
