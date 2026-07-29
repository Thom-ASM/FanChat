export function createMaterialIcon(iconName, color = 'inherit', size = '24px') {
  const icon = document.createElement('span');
  icon.classList.add('material-icons');
  icon.textContent = iconName;
  icon.style.color = color;
  icon.style.fontSize = size;
  return icon;
}
