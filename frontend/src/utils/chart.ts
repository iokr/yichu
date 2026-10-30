export const chartPalette = ['#A67C89', '#7E9A93', '#C4A484', '#8D9DB6', '#B5A8A0', '#9AA88E']

export function axisColors(dark: boolean) {
  return {
    text: dark ? '#A8A39D' : '#6B6661',
    line: dark ? '#2E2C2A' : '#EAE8E5',
    split: dark ? '#262423' : '#F4F3F1',
  }
}

export function tooltipStyle(dark: boolean) {
  return {
    backgroundColor: dark ? '#252322' : '#FFFFFF',
    borderColor: dark ? '#2E2C2A' : '#EAE8E5',
    textStyle: { color: dark ? '#F2F0EE' : '#1F1D1B', fontSize: 12 },
  }
}
