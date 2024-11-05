export const getInputClassName = (value: string, validator: (value: string) => boolean) => {
  const base = "border-2 p-2 w-full rounded outline-2 caret-blue-500 text-xl";
  if (value.length === 0) return base;
  return validator(value)
    ? `border-green-500 outline-green-500 ${base}`
    : `border-red-500 outline-red-500 ${base}`;
};
