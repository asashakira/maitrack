export type BaseEntity = {
  id: string;
  createdAt: number;
  updatedAt: number;
};

export type Entity<T> = {
  [K in keyof T]: T[K];
} & BaseEntity;

export type User = Entity<{
  gameName: string;
  tagLine: string;
  rating: string;
  seasonPlayCount: string;
  totalPlayCount: string;
}>;

export type UserRecord = Entity<{
  title: string;
  percentage: string;
  playedAt: string;
  deluxScore: string;
  maxDeluxScore: string;
  imageUrl: string;
}>;
