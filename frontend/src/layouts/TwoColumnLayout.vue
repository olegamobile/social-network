<template>
    <div class="two-column-layout">
      <aside class="sidebar">
        <slot name="sidebar"></slot>
      </aside>
      <main class="main-content">
        <slot name="main"></slot>
      </main>
    </div>
  </template>
  
  <script setup>
  // No props or reactive state needed for this simple layout
  // The setup script tag automatically exposes everything to the template
  </script>
  
  <style scoped>
  .two-column-layout {
    display: flex;
    width: 100%;
    min-height: 100vh;
  }
  
  .sidebar {
    width: 280px;
    padding: 20px;
    background-color: #f5f5f5;
    border-right: 1px solid #e0e0e0;
  }
  
  .main-content {
    flex: 1;
    padding: 20px;
    background-color: #ffffff;
  }
  
  /* Responsive behavior */
  @media (max-width: 768px) {
    .two-column-layout {
      flex-direction: column;
    }
    
    .sidebar {
      width: 100%;
      border-right: none;
      border-bottom: 1px solid #e0e0e0;
    }
  }
  </style>