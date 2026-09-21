<script lang="ts">
  import { cn } from "$lib/utils";

  let { level = "waspada", label, size = "md", pulse = true } = $props<{
    level?: "waspada" | "siaga" | "awas" | "normal" | "aman",
    label: string,
    size?: "sm" | "md" | "lg",
    pulse?: boolean
  }>();

  const VARIANTS: Record<string, any> = {
    waspada: { dot: "bg-warning", text: "text-warning", chip: "bg-warning/10 text-warning border-warning/30" },
    siaga: { dot: "bg-secondary", text: "text-secondary", chip: "bg-secondary/10 text-secondary border-secondary/30" },
    awas: { dot: "bg-destructive", text: "text-destructive", chip: "bg-destructive/10 text-destructive border-destructive/30" },
    normal: { dot: "bg-safe", text: "text-safe", chip: "bg-safe/10 text-safe border-safe/30" },
    aman: { dot: "bg-safe", text: "text-safe", chip: "bg-safe/10 text-safe border-safe/30" },
  };

  const sizes: Record<string, string> = {
    sm: "text-[0.75rem] px-2 py-0.5 gap-1.5",
    md: "text-[0.8125rem] px-2.5 py-1 gap-2",
    lg: "text-[0.875rem] px-3 py-1.5 gap-2",
  };

  let v = $derived(VARIANTS[level as string] || VARIANTS.waspada);
</script>

<span class={cn("inline-flex items-center rounded-full border font-semibold", v.chip, sizes[size])}>
  <span class="relative flex">
    <span class={cn("w-2 h-2 rounded-full", v.dot)}></span>
    {#if pulse}
      <span class={cn("absolute inset-0 w-2 h-2 rounded-full animate-ping opacity-60", v.dot)}></span>
    {/if}
  </span>
  {label}
</span>
